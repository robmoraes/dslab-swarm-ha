package config

import (
	"os"
	"strings"
	"time"
)

type BuildInfo struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
}

type AWSMetadataConfig struct {
	Mode     string        `json:"mode"`
	Timeout  time.Duration `json:"timeout"`
	CacheTTL time.Duration `json:"cache_ttl"`
}

type Config struct {
	AppName     string            `json:"app_name"`
	Environment string            `json:"environment"`
	ListenAddr  string            `json:"listen_addr"`
	Build       BuildInfo         `json:"build"`
	AWSMetadata AWSMetadataConfig `json:"aws_metadata"`
}

func Load(build BuildInfo) Config {
	port := env("PORT", "8080")
	listenAddr := env("LISTEN_ADDR", "")
	if listenAddr == "" {
		if strings.HasPrefix(port, ":") {
			listenAddr = port
		} else {
			listenAddr = ":" + port
		}
	}

	metadataMode := strings.ToLower(env("AWS_EC2_METADATA", "auto"))
	if isTruthy(os.Getenv("AWS_EC2_METADATA_DISABLED")) {
		metadataMode = "disabled"
	}
	switch metadataMode {
	case "auto", "enabled", "disabled":
	default:
		metadataMode = "auto"
	}

	return Config{
		AppName:     env("APP_NAME", "dslab-whoami"),
		Environment: env("APP_ENV", "lab"),
		ListenAddr:  listenAddr,
		Build:       build,
		AWSMetadata: AWSMetadataConfig{
			Mode:     metadataMode,
			Timeout:  durationEnv("AWS_EC2_METADATA_TIMEOUT", 250*time.Millisecond),
			CacheTTL: durationEnv("AWS_EC2_METADATA_CACHE_TTL", 30*time.Second),
		},
	}
}

func env(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func durationEnv(key string, fallback time.Duration) time.Duration {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}

	value, err := time.ParseDuration(raw)
	if err != nil || value <= 0 {
		return fallback
	}

	return value
}

func isTruthy(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "t", "true", "y", "yes", "on":
		return true
	default:
		return false
	}
}
