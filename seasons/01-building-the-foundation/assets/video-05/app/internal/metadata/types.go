package metadata

import "time"

type Snapshot struct {
	Request   RequestInfo   `json:"request"`
	App       AppInfo       `json:"app"`
	Runtime   RuntimeInfo   `json:"runtime"`
	Host      HostInfo      `json:"host"`
	Container ContainerInfo `json:"container"`
	Network   NetworkInfo   `json:"network"`
	Cloud     CloudInfo     `json:"cloud"`
	Hints     []string      `json:"hints"`
}

type RequestInfo struct {
	ID            string              `json:"id"`
	Sequence      uint64              `json:"sequence"`
	ReceivedAt    time.Time           `json:"received_at"`
	Method        string              `json:"method"`
	Path          string              `json:"path"`
	Query         string              `json:"query,omitempty"`
	Host          string              `json:"host"`
	RemoteAddress string              `json:"remote_address"`
	RemoteIP      string              `json:"remote_ip,omitempty"`
	RemotePort    string              `json:"remote_port,omitempty"`
	LocalAddress  string              `json:"local_address,omitempty"`
	Protocol      string              `json:"protocol"`
	Scheme        string              `json:"scheme"`
	UserAgent     string              `json:"user_agent,omitempty"`
	Headers       map[string]string   `json:"headers"`
	LoadBalancer  LoadBalancerHeaders `json:"load_balancer"`
}

type LoadBalancerHeaders struct {
	Forwarded       string   `json:"forwarded,omitempty"`
	ForwardedFor    []string `json:"forwarded_for,omitempty"`
	ForwardedHost   string   `json:"forwarded_host,omitempty"`
	ForwardedProto  string   `json:"forwarded_proto,omitempty"`
	ForwardedPort   string   `json:"forwarded_port,omitempty"`
	RealIP          string   `json:"real_ip,omitempty"`
	AmznTraceID     string   `json:"amzn_trace_id,omitempty"`
	Via             string   `json:"via,omitempty"`
	RequestID       string   `json:"request_id,omitempty"`
	OriginalURI     string   `json:"original_uri,omitempty"`
	OriginalForward string   `json:"original_forwarded_for,omitempty"`
}

type AppInfo struct {
	Name        string `json:"name"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
	Commit      string `json:"commit"`
	BuildDate   string `json:"build_date"`
}

type RuntimeInfo struct {
	Timestamp        time.Time `json:"timestamp"`
	Uptime           string    `json:"uptime"`
	GoVersion        string    `json:"go_version"`
	GOOS             string    `json:"goos"`
	GOARCH           string    `json:"goarch"`
	CPUs             int       `json:"cpus"`
	Goroutines       int       `json:"goroutines"`
	PID              int       `json:"pid"`
	PPID             int       `json:"ppid"`
	MemoryAllocBytes uint64    `json:"memory_alloc_bytes"`
	MemorySysBytes   uint64    `json:"memory_sys_bytes"`
	NumGC            uint32    `json:"num_gc"`
}

type HostInfo struct {
	Hostname string            `json:"hostname"`
	Env      map[string]string `json:"env,omitempty"`
}

type ContainerInfo struct {
	ID          string    `json:"id,omitempty"`
	Hostname    string    `json:"hostname"`
	CGroup      []string  `json:"cgroup,omitempty"`
	DetectedBy  string    `json:"detected_by,omitempty"`
	SwarmLabels SwarmInfo `json:"swarm"`
}

type SwarmInfo struct {
	NodeID       string `json:"node_id,omitempty"`
	NodeHostname string `json:"node_hostname,omitempty"`
	ServiceName  string `json:"service_name,omitempty"`
	TaskID       string `json:"task_id,omitempty"`
	TaskName     string `json:"task_name,omitempty"`
	TaskSlot     string `json:"task_slot,omitempty"`
}

type NetworkInfo struct {
	Interfaces []NetworkInterface `json:"interfaces"`
}

type NetworkInterface struct {
	Name         string   `json:"name"`
	HardwareAddr string   `json:"hardware_addr,omitempty"`
	Flags        []string `json:"flags,omitempty"`
	Addresses    []string `json:"addresses,omitempty"`
}

type CloudInfo struct {
	Enabled          bool      `json:"enabled"`
	Available        bool      `json:"available"`
	Provider         string    `json:"provider,omitempty"`
	InstanceID       string    `json:"instance_id,omitempty"`
	InstanceType     string    `json:"instance_type,omitempty"`
	AvailabilityZone string    `json:"availability_zone,omitempty"`
	Region           string    `json:"region,omitempty"`
	PrivateIPv4      string    `json:"private_ipv4,omitempty"`
	PublicIPv4       string    `json:"public_ipv4,omitempty"`
	CheckedAt        time.Time `json:"checked_at,omitempty"`
	CacheUntil       time.Time `json:"cache_until,omitempty"`
	Error            string    `json:"error,omitempty"`
}
