package config

import (
	"testing"
	"time"
)

func TestLoadUsesPortWhenListenAddrIsEmpty(t *testing.T) {
	t.Setenv("PORT", "9090")

	cfg := Load(BuildInfo{Version: "test"})

	if cfg.ListenAddr != ":9090" {
		t.Fatalf("expected :9090, got %q", cfg.ListenAddr)
	}
}

func TestLoadListenAddrOverridesPort(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("LISTEN_ADDR", "127.0.0.1:18080")

	cfg := Load(BuildInfo{Version: "test"})

	if cfg.ListenAddr != "127.0.0.1:18080" {
		t.Fatalf("expected explicit listen addr, got %q", cfg.ListenAddr)
	}
}

func TestLoadAWSMetadataDisabled(t *testing.T) {
	t.Setenv("AWS_EC2_METADATA", "enabled")
	t.Setenv("AWS_EC2_METADATA_DISABLED", "true")
	t.Setenv("AWS_EC2_METADATA_TIMEOUT", "750ms")

	cfg := Load(BuildInfo{Version: "test"})

	if cfg.AWSMetadata.Mode != "disabled" {
		t.Fatalf("expected disabled metadata, got %q", cfg.AWSMetadata.Mode)
	}
	if cfg.AWSMetadata.Timeout != 750*time.Millisecond {
		t.Fatalf("expected timeout 750ms, got %s", cfg.AWSMetadata.Timeout)
	}
}
