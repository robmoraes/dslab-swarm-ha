package metadata

import (
	"context"
	"net"
	"os"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/robmoraes/dslab/assets/video-05/app/internal/config"
)

var (
	containerIDPattern = regexp.MustCompile(`(?i)(?:docker[-/]|cri-containerd-|containerd[-/])?([0-9a-f]{64})(?:\.scope)?`)
	hexHostPattern     = regexp.MustCompile(`(?i)^[0-9a-f]{12,64}$`)
)

type Collector struct {
	cfg                   config.Config
	startedAt             time.Time
	hostname              string
	containerID           string
	containerIDDetectedBy string
	cgroup                []string
	cloud                 *CloudCollector
}

func NewCollector(cfg config.Config) *Collector {
	hostname, _ := os.Hostname()
	cgroup := readLines("/proc/self/cgroup")
	containerID, detectedBy := detectContainerID(hostname, cgroup)

	return &Collector{
		cfg:                   cfg,
		startedAt:             time.Now().UTC(),
		hostname:              hostname,
		containerID:           containerID,
		containerIDDetectedBy: detectedBy,
		cgroup:                cgroup,
		cloud:                 NewCloudCollector(cfg.AWSMetadata),
	}
}

func (c *Collector) Snapshot(ctx context.Context, request RequestInfo) Snapshot {
	now := time.Now().UTC()
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	return Snapshot{
		Request: request,
		App: AppInfo{
			Name:        c.cfg.AppName,
			Environment: c.cfg.Environment,
			Version:     c.cfg.Build.Version,
			Commit:      c.cfg.Build.Commit,
			BuildDate:   c.cfg.Build.Date,
		},
		Runtime: RuntimeInfo{
			Timestamp:        now,
			Uptime:           now.Sub(c.startedAt).Round(time.Millisecond).String(),
			GoVersion:        runtime.Version(),
			GOOS:             runtime.GOOS,
			GOARCH:           runtime.GOARCH,
			CPUs:             runtime.NumCPU(),
			Goroutines:       runtime.NumGoroutine(),
			PID:              os.Getpid(),
			PPID:             os.Getppid(),
			MemoryAllocBytes: mem.Alloc,
			MemorySysBytes:   mem.Sys,
			NumGC:            mem.NumGC,
		},
		Host: HostInfo{
			Hostname: c.hostname,
			Env:      selectedEnv(),
		},
		Container: ContainerInfo{
			ID:         c.containerID,
			Hostname:   c.hostname,
			CGroup:     c.cgroup,
			DetectedBy: c.containerIDDetectedBy,
			SwarmLabels: SwarmInfo{
				NodeID:       firstEnv("NODE_ID", "DOCKER_NODE_ID"),
				NodeHostname: firstEnv("NODE_HOSTNAME", "DOCKER_NODE_HOSTNAME"),
				ServiceName:  firstEnv("SERVICE_NAME", "DOCKER_SERVICE_NAME", "SWARM_SERVICE_NAME"),
				TaskID:       firstEnv("TASK_ID", "DOCKER_TASK_ID"),
				TaskName:     firstEnv("TASK_NAME", "DOCKER_TASK_NAME"),
				TaskSlot:     firstEnv("TASK_SLOT", "DOCKER_TASK_SLOT"),
			},
		},
		Network: collectNetwork(),
		Cloud:   c.cloud.Collect(ctx),
		Hints: []string{
			"Compare request.sequence, host.hostname, container.id and network.interfaces across calls to see which container answered.",
			"Behind an ALB or reverse proxy, check request.load_balancer forwarded fields and x-amzn-trace-id.",
			"In Docker Swarm, set service env templates such as NODE_HOSTNAME, SERVICE_NAME, TASK_ID and TASK_SLOT to expose scheduler placement.",
		},
	}
}

func detectContainerID(hostname string, cgroup []string) (string, string) {
	for _, line := range cgroup {
		matches := containerIDPattern.FindStringSubmatch(line)
		if len(matches) == 2 {
			return strings.ToLower(matches[1]), "/proc/self/cgroup"
		}
	}

	if hexHostPattern.MatchString(hostname) {
		return strings.ToLower(hostname), "hostname"
	}

	return "", ""
}

func collectNetwork() NetworkInfo {
	interfaces, err := net.Interfaces()
	if err != nil {
		return NetworkInfo{}
	}

	result := make([]NetworkInterface, 0, len(interfaces))
	for _, iface := range interfaces {
		item := NetworkInterface{
			Name:         iface.Name,
			HardwareAddr: iface.HardwareAddr.String(),
			Flags:        splitFlags(iface.Flags.String()),
		}

		addrs, err := iface.Addrs()
		if err == nil {
			for _, addr := range addrs {
				item.Addresses = append(item.Addresses, addr.String())
			}
			sort.Strings(item.Addresses)
		}

		result = append(result, item)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return NetworkInfo{Interfaces: result}
}

func splitFlags(raw string) []string {
	if raw == "" {
		return nil
	}
	return strings.Split(raw, "|")
}

func selectedEnv() map[string]string {
	keys := []string{
		"APP_NAME",
		"APP_ENV",
		"HOSTNAME",
		"PORT",
		"LISTEN_ADDR",
		"AWS_REGION",
		"AWS_DEFAULT_REGION",
		"AWS_EC2_METADATA",
		"AWS_EC2_METADATA_DISABLED",
		"NODE_ID",
		"NODE_HOSTNAME",
		"DOCKER_NODE_ID",
		"DOCKER_NODE_HOSTNAME",
		"SERVICE_NAME",
		"DOCKER_SERVICE_NAME",
		"SWARM_SERVICE_NAME",
		"TASK_ID",
		"TASK_NAME",
		"TASK_SLOT",
		"DOCKER_TASK_ID",
		"DOCKER_TASK_NAME",
		"DOCKER_TASK_SLOT",
	}

	values := make(map[string]string)
	for _, key := range keys {
		if value, ok := os.LookupEnv(key); ok {
			values[key] = value
		}
	}
	return values
}

func firstEnv(keys ...string) string {
	for _, key := range keys {
		if value := strings.TrimSpace(os.Getenv(key)); value != "" {
			return value
		}
	}
	return ""
}

func readLines(path string) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	rawLines := strings.Split(strings.TrimSpace(string(data)), "\n")
	lines := make([]string, 0, len(rawLines))
	for _, line := range rawLines {
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines
}
