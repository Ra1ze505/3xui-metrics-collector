package poller

import (
	"time"

	"github.com/andrejmatveev/3xui-metrics-collector/internal/xui"
)

type Snapshot struct {
	PanelName      string
	Up             bool
	ScrapeDuration time.Duration
	Errors         map[string]error
	Timestamp      time.Time

	ServerStatus *xui.ServerStatus
	Nodes        []xui.Node
	Inbounds     []xui.Inbound
	Clients      []xui.ClientWithAttachments
	OnlineEmails map[string]bool
	LastOnline   map[string]int64
	Outbounds    []xui.OutboundTraffic

	NodeNameByID map[int]string

	// cumulative per-endpoint error counts
	ErrorCounts map[string]uint64
}

func (s *Snapshot) NodeName(nodeID int) string {
	if nodeID == 0 {
		return s.PanelName
	}
	if name, ok := s.NodeNameByID[nodeID]; ok {
		return name
	}
	return "unknown"
}

func buildNodeNameByID(panelName string, nodes []xui.Node) map[int]string {
	m := map[int]string{0: panelName}
	for _, n := range nodes {
		name := n.Name
		if name == "" {
			name = n.Remark
		}
		if name == "" {
			name = n.Address
		}
		m[n.ID] = name
	}
	return m
}
