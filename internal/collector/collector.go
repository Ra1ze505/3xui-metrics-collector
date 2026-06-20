package collector

import (
	"strconv"

	"github.com/andrejmatveev/3xui-metrics-collector/internal/poller"
	"github.com/andrejmatveev/3xui-metrics-collector/internal/xui"
	"github.com/prometheus/client_golang/prometheus"
)

type SnapshotSource interface {
	Snapshot() *poller.Snapshot
}

type Collector struct {
	sources []SnapshotSource
}

func New(sources ...SnapshotSource) *Collector {
	return &Collector{sources: sources}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, d := range allDescs {
		ch <- d
	}
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	for _, src := range c.sources {
		snap := src.Snapshot()
		if snap == nil {
			continue
		}
		c.collectSnapshot(ch, snap)
	}
}

var allDescs = []*prometheus.Desc{
	// exporter
	descUp,
	descScrapeDuration,
	descScrapeErrors,
	descPanelInfo,
	// nodes
	descNodeUp,
	descNodeInfo,
	descNodeCPUPercent,
	descNodeMemPercent,
	descNodeOnlineClients,
	descNodeClientCount,
	descNodeInboundCount,
	descNodeDepletedCount,
	descNodeLatencyMs,
	descNodeLastHeartbeat,
	descNodeUptime,
	descNodeXrayUp,
	// server
	descServerCPU,
	descServerCPUCores,
	descServerMemBytes,
	descServerSwapBytes,
	descServerDiskBytes,
	descServerLoad,
	descServerNetTraffic,
	descServerNetIO,
	descServerTCP,
	descServerUDP,
	descServerUptime,
	descXrayUp,
	descXrayInfo,
	// inbounds
	descInboundUp,
	descInboundDown,
	descInboundEnabled,
	descInboundClientCount,
	descInboundExpiry,
	// clients
	descClientUp,
	descClientDown,
	descClientEnabled,
	descClientOnline,
	descClientLastOnline,
	descClientTrafficLimit,
	descClientExpiry,
	descClientInboundCount,
	descClientGroupInfo,
	descClientTrafficUp,
	descClientTrafficDown,
	// outbounds
	descOutboundUp,
	descOutboundDown,
}

var (
	descUp             = prometheus.NewDesc("xui_up", "Whether the last scrape cycle succeeded for the panel.", []string{"panel"}, nil)
	descScrapeDuration = prometheus.NewDesc("xui_scrape_duration_seconds", "Duration of the last scrape cycle.", []string{"panel"}, nil)
	descScrapeErrors   = prometheus.NewDesc("xui_scrape_errors_total", "Total scrape errors per endpoint.", []string{"panel", "endpoint"}, nil)
	descPanelInfo      = prometheus.NewDesc("xui_panel_info", "Panel version info.", []string{"panel", "version"}, nil)

	descNodeUp            = prometheus.NewDesc("xui_node_up", "Whether the node is up.", []string{"panel", "node", "role"}, nil)
	descNodeInfo          = prometheus.NewDesc("xui_node_info", "Node info.", []string{"panel", "node", "role", "panel_version", "xray_version", "address"}, nil)
	descNodeCPUPercent    = prometheus.NewDesc("xui_node_cpu_percent", "Node CPU usage percent.", []string{"panel", "node", "role"}, nil)
	descNodeMemPercent    = prometheus.NewDesc("xui_node_mem_percent", "Node memory usage percent.", []string{"panel", "node", "role"}, nil)
	descNodeOnlineClients = prometheus.NewDesc("xui_node_online_clients", "Online clients on the node.", []string{"panel", "node", "role"}, nil)
	descNodeClientCount   = prometheus.NewDesc("xui_node_client_count", "Total clients on the node.", []string{"panel", "node", "role"}, nil)
	descNodeInboundCount  = prometheus.NewDesc("xui_node_inbound_count", "Inbound count on the node.", []string{"panel", "node", "role"}, nil)
	descNodeDepletedCount = prometheus.NewDesc("xui_node_depleted_count", "Depleted clients on the node.", []string{"panel", "node", "role"}, nil)
	descNodeLatencyMs     = prometheus.NewDesc("xui_node_latency_ms", "Node latency in milliseconds.", []string{"panel", "node", "role"}, nil)
	descNodeLastHeartbeat = prometheus.NewDesc("xui_node_last_heartbeat_seconds", "Last heartbeat unix timestamp.", []string{"panel", "node", "role"}, nil)
	descNodeUptime        = prometheus.NewDesc("xui_node_uptime_seconds", "Node uptime in seconds.", []string{"panel", "node", "role"}, nil)
	descNodeXrayUp        = prometheus.NewDesc("xui_node_xray_up", "Whether xray is running on the node.", []string{"panel", "node", "role"}, nil)

	descServerCPU        = prometheus.NewDesc("xui_server_cpu_percent", "Master server CPU percent.", []string{"panel"}, nil)
	descServerCPUCores   = prometheus.NewDesc("xui_server_cpu_cores", "Master server CPU cores.", []string{"panel"}, nil)
	descServerMemBytes   = prometheus.NewDesc("xui_server_mem_bytes", "Master server memory bytes.", []string{"panel", "state"}, nil)
	descServerSwapBytes  = prometheus.NewDesc("xui_server_swap_bytes", "Master server swap bytes.", []string{"panel", "state"}, nil)
	descServerDiskBytes  = prometheus.NewDesc("xui_server_disk_bytes", "Master server disk bytes.", []string{"panel", "state"}, nil)
	descServerLoad       = prometheus.NewDesc("xui_server_load", "Master server load average.", []string{"panel", "period"}, nil)
	descServerNetTraffic = prometheus.NewDesc("xui_server_net_traffic_bytes_total", "Master server cumulative network traffic.", []string{"panel", "direction"}, nil)
	descServerNetIO      = prometheus.NewDesc("xui_server_net_io_bytes", "Master server current network IO rate.", []string{"panel", "direction"}, nil)
	descServerTCP        = prometheus.NewDesc("xui_server_tcp_connections", "Master server TCP connections.", []string{"panel"}, nil)
	descServerUDP        = prometheus.NewDesc("xui_server_udp_connections", "Master server UDP connections.", []string{"panel"}, nil)
	descServerUptime     = prometheus.NewDesc("xui_server_uptime_seconds", "Master server uptime.", []string{"panel"}, nil)
	descXrayUp           = prometheus.NewDesc("xui_xray_up", "Whether xray is running on master.", []string{"panel"}, nil)
	descXrayInfo         = prometheus.NewDesc("xui_xray_info", "Xray version on master.", []string{"panel", "version"}, nil)

	descInboundUp          = prometheus.NewDesc("xui_inbound_up_bytes_total", "Inbound upload bytes.", []string{"panel", "node", "inbound_id", "remark", "tag", "protocol", "port"}, nil)
	descInboundDown        = prometheus.NewDesc("xui_inbound_down_bytes_total", "Inbound download bytes.", []string{"panel", "node", "inbound_id", "remark", "tag", "protocol", "port"}, nil)
	descInboundEnabled     = prometheus.NewDesc("xui_inbound_enabled", "Whether inbound is enabled.", []string{"panel", "node", "inbound_id", "remark", "tag", "protocol", "port"}, nil)
	descInboundClientCount = prometheus.NewDesc("xui_inbound_client_count", "Clients on inbound.", []string{"panel", "node", "inbound_id", "remark", "tag", "protocol", "port"}, nil)
	descInboundExpiry      = prometheus.NewDesc("xui_inbound_expiry_timestamp_seconds", "Inbound expiry timestamp.", []string{"panel", "node", "inbound_id", "remark", "tag", "protocol", "port"}, nil)

	descClientUp           = prometheus.NewDesc("xui_client_up_bytes_total", "Client upload bytes per inbound.", []string{"panel", "node", "email", "inbound_id", "inbound_remark", "protocol"}, nil)
	descClientDown         = prometheus.NewDesc("xui_client_down_bytes_total", "Client download bytes per inbound.", []string{"panel", "node", "email", "inbound_id", "inbound_remark", "protocol"}, nil)
	descClientEnabled      = prometheus.NewDesc("xui_client_enabled", "Whether client is enabled.", []string{"panel", "email"}, nil)
	descClientOnline       = prometheus.NewDesc("xui_client_online", "Whether client is online.", []string{"panel", "email"}, nil)
	descClientLastOnline   = prometheus.NewDesc("xui_client_last_online_timestamp_seconds", "Client last online timestamp.", []string{"panel", "email"}, nil)
	descClientTrafficLimit = prometheus.NewDesc("xui_client_traffic_limit_bytes", "Client traffic limit in bytes (0=unlimited).", []string{"panel", "email"}, nil)
	descClientExpiry       = prometheus.NewDesc("xui_client_expiry_timestamp_seconds", "Client expiry timestamp.", []string{"panel", "email"}, nil)
	descClientInboundCount = prometheus.NewDesc("xui_client_inbound_count", "Number of inbounds attached to client.", []string{"panel", "email"}, nil)
	descClientGroupInfo    = prometheus.NewDesc("xui_client_group_info", "Client group info.", []string{"panel", "email", "group"}, nil)

	descClientTrafficUp   = prometheus.NewDesc("xui_client_traffic_up_bytes_total", "Authoritative per-client cumulative upload bytes, email-keyed (aggregated across all nodes by 3x-ui).", []string{"panel", "email"}, nil)
	descClientTrafficDown = prometheus.NewDesc("xui_client_traffic_down_bytes_total", "Authoritative per-client cumulative download bytes, email-keyed (aggregated across all nodes by 3x-ui).", []string{"panel", "email"}, nil)

	descOutboundUp   = prometheus.NewDesc("xui_outbound_up_bytes_total", "Outbound upload bytes.", []string{"panel", "tag"}, nil)
	descOutboundDown = prometheus.NewDesc("xui_outbound_down_bytes_total", "Outbound download bytes.", []string{"panel", "tag"}, nil)
)

func (c *Collector) collectSnapshot(ch chan<- prometheus.Metric, snap *poller.Snapshot) {
	panel := snap.PanelName

	upVal := 0.0
	if snap.Up {
		upVal = 1
	}
	ch <- prometheus.MustNewConstMetric(descUp, prometheus.GaugeValue, upVal, panel)
	ch <- prometheus.MustNewConstMetric(descScrapeDuration, prometheus.GaugeValue, snap.ScrapeDuration.Seconds(), panel)

	for endpoint, count := range snap.ErrorCounts {
		ch <- prometheus.MustNewConstMetric(descScrapeErrors, prometheus.CounterValue, float64(count), panel, endpoint)
	}

	if snap.ServerStatus != nil {
		version := snap.ServerStatus.PanelVersion
		if version != "" {
			ch <- prometheus.MustNewConstMetric(descPanelInfo, prometheus.GaugeValue, 1, panel, version)
		}
		c.collectMasterNode(ch, snap)
		c.collectServer(ch, snap)
	}

	for _, node := range snap.Nodes {
		c.collectRemoteNode(ch, snap, node)
	}

	c.collectInbounds(ch, snap)
	c.collectClients(ch, snap)

	for _, ob := range snap.Outbounds {
		ch <- prometheus.MustNewConstMetric(descOutboundUp, prometheus.CounterValue, float64(ob.Up), panel, ob.Tag)
		ch <- prometheus.MustNewConstMetric(descOutboundDown, prometheus.CounterValue, float64(ob.Down), panel, ob.Tag)
	}
}

func (c *Collector) collectMasterNode(ch chan<- prometheus.Metric, snap *poller.Snapshot) {
	panel := snap.PanelName
	node := snap.PanelName
	role := "master"
	status := snap.ServerStatus

	xrayUp := 0.0
	if status.Xray.State == "running" {
		xrayUp = 1
	}
	masterUp := xrayUp

	memPct := 0.0
	if status.Mem.Total > 0 {
		memPct = float64(status.Mem.Current) / float64(status.Mem.Total) * 100
	}

	onlineCount := float64(len(snap.OnlineEmails))
	clientCount := float64(len(snap.Clients))
	localInbounds := 0
	for _, ib := range snap.Inbounds {
		if ib.NodeID == 0 {
			localInbounds++
		}
	}

	address := status.PublicIP.IPv4
	if address == "" {
		address = status.PublicIP.IPv6
	}

	ch <- prometheus.MustNewConstMetric(descNodeUp, prometheus.GaugeValue, masterUp, panel, node, role)
	ch <- prometheus.MustNewConstMetric(descNodeInfo, prometheus.GaugeValue, 1, panel, node, role, status.PanelVersion, status.Xray.Version, address)
	ch <- prometheus.MustNewConstMetric(descNodeCPUPercent, prometheus.GaugeValue, status.CPU, panel, node, role)
	ch <- prometheus.MustNewConstMetric(descNodeMemPercent, prometheus.GaugeValue, memPct, panel, node, role)
	ch <- prometheus.MustNewConstMetric(descNodeOnlineClients, prometheus.GaugeValue, onlineCount, panel, node, role)
	ch <- prometheus.MustNewConstMetric(descNodeClientCount, prometheus.GaugeValue, clientCount, panel, node, role)
	ch <- prometheus.MustNewConstMetric(descNodeInboundCount, prometheus.GaugeValue, float64(localInbounds), panel, node, role)
	ch <- prometheus.MustNewConstMetric(descNodeDepletedCount, prometheus.GaugeValue, 0, panel, node, role)
	ch <- prometheus.MustNewConstMetric(descNodeUptime, prometheus.GaugeValue, float64(status.Uptime), panel, node, role)
	ch <- prometheus.MustNewConstMetric(descNodeXrayUp, prometheus.GaugeValue, xrayUp, panel, node, role)
}

func (c *Collector) collectRemoteNode(ch chan<- prometheus.Metric, snap *poller.Snapshot, node xui.Node) {
	panel := snap.PanelName
	nodeName := snap.NodeName(node.ID)
	role := "node"

	nodeUp := 0.0
	if node.Enable && node.Status != "offline" && node.Status != "down" {
		nodeUp = 1
	}
	xrayUp := 0.0
	if node.XrayState == "running" {
		xrayUp = 1
	}

	address := node.Address
	if node.Port > 0 {
		address = address + ":" + strconv.Itoa(node.Port)
	}

	ch <- prometheus.MustNewConstMetric(descNodeUp, prometheus.GaugeValue, nodeUp, panel, nodeName, role)
	ch <- prometheus.MustNewConstMetric(descNodeInfo, prometheus.GaugeValue, 1, panel, nodeName, role, node.PanelVersion, node.XrayVersion, address)
	ch <- prometheus.MustNewConstMetric(descNodeCPUPercent, prometheus.GaugeValue, node.CPUPct, panel, nodeName, role)
	ch <- prometheus.MustNewConstMetric(descNodeMemPercent, prometheus.GaugeValue, node.MemPct, panel, nodeName, role)
	ch <- prometheus.MustNewConstMetric(descNodeOnlineClients, prometheus.GaugeValue, float64(node.OnlineCount), panel, nodeName, role)
	ch <- prometheus.MustNewConstMetric(descNodeClientCount, prometheus.GaugeValue, float64(node.ClientCount), panel, nodeName, role)
	ch <- prometheus.MustNewConstMetric(descNodeInboundCount, prometheus.GaugeValue, float64(node.InboundCount), panel, nodeName, role)
	ch <- prometheus.MustNewConstMetric(descNodeDepletedCount, prometheus.GaugeValue, float64(node.DepletedCount), panel, nodeName, role)
	ch <- prometheus.MustNewConstMetric(descNodeLatencyMs, prometheus.GaugeValue, float64(node.LatencyMs), panel, nodeName, role)
	if node.LastHeartbeat > 0 {
		ch <- prometheus.MustNewConstMetric(descNodeLastHeartbeat, prometheus.GaugeValue, float64(node.LastHeartbeat), panel, nodeName, role)
	}
	ch <- prometheus.MustNewConstMetric(descNodeUptime, prometheus.GaugeValue, float64(node.UptimeSecs), panel, nodeName, role)
	ch <- prometheus.MustNewConstMetric(descNodeXrayUp, prometheus.GaugeValue, xrayUp, panel, nodeName, role)
}

func (c *Collector) collectServer(ch chan<- prometheus.Metric, snap *poller.Snapshot) {
	panel := snap.PanelName
	s := snap.ServerStatus

	ch <- prometheus.MustNewConstMetric(descServerCPU, prometheus.GaugeValue, s.CPU, panel)
	ch <- prometheus.MustNewConstMetric(descServerCPUCores, prometheus.GaugeValue, float64(s.CPUCores), panel)
	ch <- prometheus.MustNewConstMetric(descServerMemBytes, prometheus.GaugeValue, float64(s.Mem.Current), panel, "used")
	ch <- prometheus.MustNewConstMetric(descServerMemBytes, prometheus.GaugeValue, float64(s.Mem.Total), panel, "total")
	ch <- prometheus.MustNewConstMetric(descServerSwapBytes, prometheus.GaugeValue, float64(s.Swap.Current), panel, "used")
	ch <- prometheus.MustNewConstMetric(descServerSwapBytes, prometheus.GaugeValue, float64(s.Swap.Total), panel, "total")
	ch <- prometheus.MustNewConstMetric(descServerDiskBytes, prometheus.GaugeValue, float64(s.Disk.Current), panel, "used")
	ch <- prometheus.MustNewConstMetric(descServerDiskBytes, prometheus.GaugeValue, float64(s.Disk.Total), panel, "total")

	periods := []string{"1m", "5m", "15m"}
	for i, load := range s.Loads {
		if i >= len(periods) {
			break
		}
		ch <- prometheus.MustNewConstMetric(descServerLoad, prometheus.GaugeValue, load, panel, periods[i])
	}

	ch <- prometheus.MustNewConstMetric(descServerNetTraffic, prometheus.CounterValue, float64(s.NetTraffic.Sent), panel, "sent")
	ch <- prometheus.MustNewConstMetric(descServerNetTraffic, prometheus.CounterValue, float64(s.NetTraffic.Recv), panel, "recv")
	ch <- prometheus.MustNewConstMetric(descServerNetIO, prometheus.GaugeValue, float64(s.NetIO.Up), panel, "up")
	ch <- prometheus.MustNewConstMetric(descServerNetIO, prometheus.GaugeValue, float64(s.NetIO.Down), panel, "down")
	ch <- prometheus.MustNewConstMetric(descServerTCP, prometheus.GaugeValue, float64(s.TCPCount), panel)
	ch <- prometheus.MustNewConstMetric(descServerUDP, prometheus.GaugeValue, float64(s.UDPCount), panel)
	ch <- prometheus.MustNewConstMetric(descServerUptime, prometheus.GaugeValue, float64(s.Uptime), panel)

	xrayUp := 0.0
	if s.Xray.State == "running" {
		xrayUp = 1
	}
	ch <- prometheus.MustNewConstMetric(descXrayUp, prometheus.GaugeValue, xrayUp, panel)
	if s.Xray.Version != "" {
		ch <- prometheus.MustNewConstMetric(descXrayInfo, prometheus.GaugeValue, 1, panel, s.Xray.Version)
	}
}

func inboundLabels(snap *poller.Snapshot, ib xui.Inbound) []string {
	return []string{
		snap.PanelName,
		snap.NodeName(ib.NodeID),
		strconv.Itoa(ib.ID),
		ib.Remark,
		ib.Tag,
		ib.Protocol,
		strconv.Itoa(ib.Port),
	}
}

func (c *Collector) collectInbounds(ch chan<- prometheus.Metric, snap *poller.Snapshot) {
	for _, ib := range snap.Inbounds {
		labels := inboundLabels(snap, ib)
		enabled := 0.0
		if ib.Enable {
			enabled = 1
		}
		ch <- prometheus.MustNewConstMetric(descInboundUp, prometheus.CounterValue, float64(ib.Up), labels...)
		ch <- prometheus.MustNewConstMetric(descInboundDown, prometheus.CounterValue, float64(ib.Down), labels...)
		ch <- prometheus.MustNewConstMetric(descInboundEnabled, prometheus.GaugeValue, enabled, labels...)
		ch <- prometheus.MustNewConstMetric(descInboundClientCount, prometheus.GaugeValue, float64(len(ib.ClientStats)), labels...)
		if ib.ExpiryTime > 0 {
			ch <- prometheus.MustNewConstMetric(descInboundExpiry, prometheus.GaugeValue, float64(ib.ExpiryTime)/1000, labels...)
		}
	}
}

func (c *Collector) collectClients(ch chan<- prometheus.Metric, snap *poller.Snapshot) {
	panel := snap.PanelName

	for _, ib := range snap.Inbounds {
		node := snap.NodeName(ib.NodeID)
		for _, cs := range ib.ClientStats {
			if cs.Email == "" {
				continue
			}
			labels := []string{
				panel,
				node,
				cs.Email,
				strconv.Itoa(ib.ID),
				ib.Remark,
				ib.Protocol,
			}
			ch <- prometheus.MustNewConstMetric(descClientUp, prometheus.CounterValue, float64(cs.Up), labels...)
			ch <- prometheus.MustNewConstMetric(descClientDown, prometheus.CounterValue, float64(cs.Down), labels...)
		}
	}

	for _, cl := range snap.Clients {
		if cl.Email == "" {
			continue
		}
		enabled := 0.0
		if cl.Enable {
			enabled = 1
		}
		online := 0.0
		if snap.OnlineEmails[cl.Email] {
			online = 1
		}
		ch <- prometheus.MustNewConstMetric(descClientEnabled, prometheus.GaugeValue, enabled, panel, cl.Email)
		ch <- prometheus.MustNewConstMetric(descClientOnline, prometheus.GaugeValue, online, panel, cl.Email)
		if ts, ok := snap.LastOnline[cl.Email]; ok && ts > 0 {
			ch <- prometheus.MustNewConstMetric(descClientLastOnline, prometheus.GaugeValue, float64(ts), panel, cl.Email)
		}
		ch <- prometheus.MustNewConstMetric(descClientTrafficLimit, prometheus.GaugeValue, float64(cl.TotalGB), panel, cl.Email)
		if cl.ExpiryTime > 0 {
			ch <- prometheus.MustNewConstMetric(descClientExpiry, prometheus.GaugeValue, float64(cl.ExpiryTime)/1000, panel, cl.Email)
		}
		ch <- prometheus.MustNewConstMetric(descClientInboundCount, prometheus.GaugeValue, float64(len(cl.InboundIDs)), panel, cl.Email)
		group := cl.Group
		if group == "" {
			group = "default"
		}
		ch <- prometheus.MustNewConstMetric(descClientGroupInfo, prometheus.GaugeValue, 1, panel, cl.Email, group)

		// Authoritative per-client totals come from the email-keyed client_traffics
		// row (clients/list). Per-inbound clientStats are replicated across the
		// client's inbounds in multi-node setups and cannot be summed reliably.
		if cl.Traffic != nil {
			ch <- prometheus.MustNewConstMetric(descClientTrafficUp, prometheus.CounterValue, float64(cl.Traffic.Up), panel, cl.Email)
			ch <- prometheus.MustNewConstMetric(descClientTrafficDown, prometheus.CounterValue, float64(cl.Traffic.Down), panel, cl.Email)
		}
	}
}
