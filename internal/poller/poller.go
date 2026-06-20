package poller

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/andrejmatveev/3xui-metrics-collector/internal/xui"
)

const (
	endpointServerStatus = "server/status"
	endpointNodes        = "nodes/list"
	endpointInbounds     = "inbounds/list"
	endpointClients      = "clients/list"
	endpointOnlines      = "clients/onlines"
	endpointLastOnline   = "clients/lastOnline"
	endpointOutbounds    = "xray/getOutboundsTraffic"
)

type Poller struct {
	name             string
	client           *xui.Client
	interval         time.Duration
	requestTimeout   time.Duration
	collectOutbounds bool

	mu          sync.RWMutex
	current     *Snapshot
	errorCounts map[string]uint64
}

func New(name string, client *xui.Client, interval, requestTimeout time.Duration, collectOutbounds bool) *Poller {
	return &Poller{
		name:             name,
		client:           client,
		interval:         interval,
		requestTimeout:   requestTimeout,
		collectOutbounds: collectOutbounds,
		errorCounts:      make(map[string]uint64),
		current: &Snapshot{
			PanelName:    name,
			Errors:       make(map[string]error),
			OnlineEmails: make(map[string]bool),
			LastOnline:   make(map[string]int64),
			NodeNameByID: map[int]string{0: name},
			ErrorCounts:  make(map[string]uint64),
		},
	}
}

func (p *Poller) Name() string {
	return p.name
}

func (p *Poller) Snapshot() *Snapshot {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if p.current == nil {
		return nil
	}
	return cloneSnapshot(p.current)
}

func (p *Poller) Start(ctx context.Context) {
	p.pollOnce(ctx)

	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.pollOnce(ctx)
		}
	}
}

func (p *Poller) pollOnce(ctx context.Context) {
	start := time.Now()

	p.mu.RLock()
	prev := p.current
	p.mu.RUnlock()

	snap := &Snapshot{
		PanelName:    p.name,
		Errors:       make(map[string]error),
		Timestamp:    start,
		OnlineEmails: make(map[string]bool),
		LastOnline:   make(map[string]int64),
		ErrorCounts:  copyErrorCounts(p.errorCounts),
	}

	if prev != nil {
		snap.ServerStatus = prev.ServerStatus
		snap.Nodes = prev.Nodes
		snap.Inbounds = prev.Inbounds
		snap.Clients = prev.Clients
		snap.Outbounds = prev.Outbounds
		if prev.OnlineEmails != nil {
			for k, v := range prev.OnlineEmails {
				snap.OnlineEmails[k] = v
			}
		}
		if prev.LastOnline != nil {
			for k, v := range prev.LastOnline {
				snap.LastOnline[k] = v
			}
		}
	}

	reqCtx, cancel := context.WithTimeout(ctx, p.requestTimeout)
	defer cancel()

	var wg sync.WaitGroup
	var mu sync.Mutex

	fetch := func(endpoint string, fn func(context.Context) error) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := fn(reqCtx); err != nil {
				log.Printf("panel %s: %s: %v", p.name, endpoint, err)
				mu.Lock()
				snap.Errors[endpoint] = err
				p.errorCounts[endpoint]++
				snap.ErrorCounts[endpoint] = p.errorCounts[endpoint]
				mu.Unlock()
			}
		}()
	}

	fetch(endpointServerStatus, func(c context.Context) error {
		status, err := p.client.GetServerStatus(c)
		if err != nil {
			return err
		}
		mu.Lock()
		snap.ServerStatus = status
		mu.Unlock()
		return nil
	})

	fetch(endpointNodes, func(c context.Context) error {
		nodes, err := p.client.GetNodes(c)
		if err != nil {
			return err
		}
		mu.Lock()
		snap.Nodes = nodes
		mu.Unlock()
		return nil
	})

	fetch(endpointInbounds, func(c context.Context) error {
		inbounds, err := p.client.GetInbounds(c)
		if err != nil {
			return err
		}
		mu.Lock()
		snap.Inbounds = inbounds
		mu.Unlock()
		return nil
	})

	fetch(endpointClients, func(c context.Context) error {
		clients, err := p.client.GetClients(c)
		if err != nil {
			return err
		}
		mu.Lock()
		snap.Clients = clients
		mu.Unlock()
		return nil
	})

	fetch(endpointOnlines, func(c context.Context) error {
		onlines, err := p.client.GetOnlines(c)
		if err != nil {
			return err
		}
		mu.Lock()
		for _, email := range onlines {
			if email != "" {
				snap.OnlineEmails[email] = true
			}
		}
		mu.Unlock()
		return nil
	})

	fetch(endpointLastOnline, func(c context.Context) error {
		lastOnline, err := p.client.GetLastOnline(c)
		if err != nil {
			return err
		}
		mu.Lock()
		snap.LastOnline = lastOnline
		mu.Unlock()
		return nil
	})

	if p.collectOutbounds {
		fetch(endpointOutbounds, func(c context.Context) error {
			outbounds, err := p.client.GetOutboundsTraffic(c)
			if err != nil {
				return err
			}
			mu.Lock()
			snap.Outbounds = outbounds
			mu.Unlock()
			return nil
		})
	}

	wg.Wait()

	snap.NodeNameByID = buildNodeNameByID(p.name, snap.Nodes)
	snap.ScrapeDuration = time.Since(start)

	_, serverFailed := snap.Errors[endpointServerStatus]
	_, inboundsFailed := snap.Errors[endpointInbounds]
	snap.Up = !serverFailed && !inboundsFailed

	p.mu.Lock()
	p.current = snap
	p.mu.Unlock()
}

func copyErrorCounts(src map[string]uint64) map[string]uint64 {
	dst := make(map[string]uint64, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func cloneSnapshot(s *Snapshot) *Snapshot {
	cp := *s
	if s.Errors != nil {
		cp.Errors = make(map[string]error, len(s.Errors))
		for k, v := range s.Errors {
			cp.Errors[k] = v
		}
	}
	if s.OnlineEmails != nil {
		cp.OnlineEmails = make(map[string]bool, len(s.OnlineEmails))
		for k, v := range s.OnlineEmails {
			cp.OnlineEmails[k] = v
		}
	}
	if s.LastOnline != nil {
		cp.LastOnline = make(map[string]int64, len(s.LastOnline))
		for k, v := range s.LastOnline {
			cp.LastOnline[k] = v
		}
	}
	if s.NodeNameByID != nil {
		cp.NodeNameByID = make(map[int]string, len(s.NodeNameByID))
		for k, v := range s.NodeNameByID {
			cp.NodeNameByID[k] = v
		}
	}
	if s.ErrorCounts != nil {
		cp.ErrorCounts = make(map[string]uint64, len(s.ErrorCounts))
		for k, v := range s.ErrorCounts {
			cp.ErrorCounts[k] = v
		}
	}
	return &cp
}
