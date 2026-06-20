package xui

type MemStats struct {
	Current uint64 `json:"current"`
	Total   uint64 `json:"total"`
}

type DiskIOStats struct {
	Read  uint64 `json:"read"`
	Write uint64 `json:"write"`
}

type XrayStatus struct {
	State    string `json:"state"`
	ErrorMsg string `json:"errorMsg"`
	Version  string `json:"version"`
}

type PublicIP struct {
	IPv4 string `json:"ipv4"`
	IPv6 string `json:"ipv6"`
}

type NetIO struct {
	Up      uint64 `json:"up"`
	Down    uint64 `json:"down"`
	PktUp   uint64 `json:"pktUp"`
	PktDown uint64 `json:"pktDown"`
}

type NetTraffic struct {
	Sent    uint64 `json:"sent"`
	Recv    uint64 `json:"recv"`
	PktSent uint64 `json:"pktSent"`
	PktRecv uint64 `json:"pktRecv"`
}

type AppStats struct {
	Threads int    `json:"threads"`
	Mem     uint64 `json:"mem"`
	Uptime  uint64 `json:"uptime"`
}

type ServerStatus struct {
	CPU          float64     `json:"cpu"`
	CPUCores     int         `json:"cpuCores"`
	LogicalPro   int         `json:"logicalPro"`
	CPUSpeedMhz  float64     `json:"cpuSpeedMhz"`
	Mem          MemStats    `json:"mem"`
	Swap         MemStats    `json:"swap"`
	Disk         MemStats    `json:"disk"`
	DiskIO       DiskIOStats `json:"diskIO"`
	DiskTraffic  DiskIOStats `json:"diskTraffic"`
	Xray         XrayStatus  `json:"xray"`
	PanelVersion string      `json:"panelVersion"`
	PanelGuid    string      `json:"panelGuid"`
	Uptime       uint64      `json:"uptime"`
	Loads        []float64   `json:"loads"`
	TCPCount     int         `json:"tcpCount"`
	UDPCount     int         `json:"udpCount"`
	NetIO        NetIO       `json:"netIO"`
	NetTraffic   NetTraffic  `json:"netTraffic"`
	PublicIP     PublicIP    `json:"publicIP"`
	AppStats     AppStats    `json:"appStats"`
}

type Node struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Remark        string  `json:"remark"`
	Address       string  `json:"address"`
	Port          int     `json:"port"`
	Scheme        string  `json:"scheme"`
	BasePath      string  `json:"basePath"`
	Enable        bool    `json:"enable"`
	Status        string  `json:"status"`
	Guid          string  `json:"guid"`
	ParentGuid    string  `json:"parentGuid"`
	Transitive    bool    `json:"transitive"`
	CPUPct        float64 `json:"cpuPct"`
	MemPct        float64 `json:"memPct"`
	OnlineCount   int     `json:"onlineCount"`
	ClientCount   int     `json:"clientCount"`
	InboundCount  int     `json:"inboundCount"`
	DepletedCount int     `json:"depletedCount"`
	LatencyMs     int     `json:"latencyMs"`
	LastHeartbeat int     `json:"lastHeartbeat"`
	UptimeSecs    int     `json:"uptimeSecs"`
	XrayState     string  `json:"xrayState"`
	XrayError     string  `json:"xrayError"`
	XrayVersion   string  `json:"xrayVersion"`
	PanelVersion  string  `json:"panelVersion"`
	LastError     string  `json:"lastError"`
	ConfigDirty   bool    `json:"configDirty"`
}

type ClientTraffic struct {
	ID         int    `json:"id"`
	InboundID  int    `json:"inboundId"`
	Enable     bool   `json:"enable"`
	Email      string `json:"email"`
	Up         int64  `json:"up"`
	Down       int64  `json:"down"`
	Total      int64  `json:"total"`
	ExpiryTime int64  `json:"expiryTime"`
	Reset      int    `json:"reset"`
	LastOnline int64  `json:"lastOnline"`
	SubID      string `json:"subId"`
	UUID       string `json:"uuid"`
}

type Inbound struct {
	ID             int             `json:"id"`
	Up             int64           `json:"up"`
	Down           int64           `json:"down"`
	Total          int64           `json:"total"`
	Remark         string          `json:"remark"`
	Enable         bool            `json:"enable"`
	ExpiryTime     int64           `json:"expiryTime"`
	Listen         string          `json:"listen"`
	Port           int             `json:"port"`
	Protocol       string          `json:"protocol"`
	Tag            string          `json:"tag"`
	NodeID         int             `json:"nodeId"`
	OriginNodeGuid string          `json:"originNodeGuid"`
	ClientStats    []ClientTraffic `json:"clientStats"`
}

type ClientRecord struct {
	Email      string `json:"email"`
	Enable     bool   `json:"enable"`
	ExpiryTime int64  `json:"expiryTime"`
	TotalGB    int64  `json:"totalGB"`
	Group      string `json:"group"`
	SubID      string `json:"subId"`
	ID         int    `json:"id"`
	UUID       string `json:"uuid"`
	LimitIP    int    `json:"limitIp"`
	Reset      int    `json:"reset"`
}

type ClientWithAttachments struct {
	ClientRecord
	InboundIDs []int          `json:"inboundIds"`
	Traffic    *ClientTraffic `json:"traffic"`
}

type OutboundTraffic struct {
	Tag   string `json:"tag"`
	Up    int64  `json:"up"`
	Down  int64  `json:"down"`
	Total int64  `json:"total"`
}
