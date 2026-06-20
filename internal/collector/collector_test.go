package collector_test

import (
	"strings"
	"testing"

	"github.com/andrejmatveev/3xui-metrics-collector/internal/collector"
	"github.com/andrejmatveev/3xui-metrics-collector/internal/poller"
	"github.com/andrejmatveev/3xui-metrics-collector/internal/xui"
	"github.com/prometheus/client_golang/prometheus"
)

func TestCollectorEmitsExpectedMetrics(t *testing.T) {
	snap := &poller.Snapshot{
		PanelName:    "test-panel",
		Up:           true,
		NodeNameByID: map[int]string{0: "test-panel"},
		OnlineEmails: map[string]bool{"user@example.com": true},
		LastOnline:   map[string]int64{"user@example.com": 1710000000},
		ErrorCounts:  map[string]uint64{"clients/list": 1},
		ServerStatus: &xui.ServerStatus{
			CPU:          10,
			CPUCores:     4,
			PanelVersion: "3.3.1",
			Uptime:       100,
			Mem:          xui.MemStats{Current: 512, Total: 1024},
			Xray:         xui.XrayStatus{State: "running", Version: "25.1.1"},
			NetTraffic:   xui.NetTraffic{Sent: 1000, Recv: 2000},
			Loads:        []float64{0.1, 0.2, 0.3},
		},
		Nodes: []xui.Node{{
			ID:            1,
			Name:          "node-1",
			Enable:        true,
			Status:        "online",
			CPUPct:        20,
			MemPct:        30,
			OnlineCount:   2,
			ClientCount:   5,
			XrayState:     "running",
			LastHeartbeat: 1710000100,
		}},
		Inbounds: []xui.Inbound{{
			ID:       10,
			Up:       100,
			Down:     200,
			Remark:   "main",
			Enable:   true,
			Port:     443,
			Protocol: "vless",
			NodeID:   0,
			ClientStats: []xui.ClientTraffic{{
				Email: "user@example.com",
				Up:    50,
				Down:  75,
			}},
		}},
		Clients: []xui.ClientWithAttachments{{
			ClientRecord: xui.ClientRecord{
				Email:      "user@example.com",
				Enable:     true,
				ExpiryTime: 1893456000000,
				TotalGB:    10737418240,
				Group:      "vip",
			},
			InboundIDs: []int{10},
			Traffic:    &xui.ClientTraffic{Email: "user@example.com", Up: 12345, Down: 67890},
		}},
		Outbounds: []xui.OutboundTraffic{{Tag: "direct", Up: 1, Down: 2}},
	}

	col := collector.New(&fakePoller{name: "test-panel", snap: snap})

	ch := make(chan prometheus.Metric, 256)
	col.Collect(ch)
	close(ch)

	names := map[string]int{}
	for m := range ch {
		desc := m.Desc().String()
		for _, required := range []string{
			"xui_up", "xui_scrape_errors_total", "xui_panel_info", "xui_node_up",
			"xui_server_cpu_percent", "xui_inbound_up_bytes_total", "xui_client_up_bytes_total",
			"xui_client_online", "xui_client_group_info", "xui_outbound_up_bytes_total",
			"xui_client_traffic_up_bytes_total", "xui_client_traffic_down_bytes_total",
		} {
			if strings.Contains(desc, `fqName: "`+required+`"`) {
				names[required]++
			}
		}
	}

	required := []string{
		"xui_up",
		"xui_scrape_errors_total",
		"xui_panel_info",
		"xui_node_up",
		"xui_server_cpu_percent",
		"xui_inbound_up_bytes_total",
		"xui_client_up_bytes_total",
		"xui_client_online",
		"xui_client_group_info",
		"xui_outbound_up_bytes_total",
		"xui_client_traffic_up_bytes_total",
		"xui_client_traffic_down_bytes_total",
	}
	for _, name := range required {
		if names[name] == 0 {
			t.Fatalf("missing metric %s; got counts: %v", name, names)
		}
	}
}

type fakePoller struct {
	name string
	snap *poller.Snapshot
}

func (f *fakePoller) Snapshot() *poller.Snapshot {
	return f.snap
}
