package httpapi

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/robmoraes/dslab/assets/video-05/app/internal/config"
	"github.com/robmoraes/dslab/assets/video-05/app/internal/metadata"
)

func TestHeadersEndpointRedactsSensitiveHeaders(t *testing.T) {
	cfg := config.Load(config.BuildInfo{Version: "test"})
	cfg.AWSMetadata.Mode = "disabled"

	server := NewServer(cfg, metadata.NewCollector(cfg), log.New(testWriter{t: t}, "", 0))

	req := httptest.NewRequest(http.MethodGet, "/headers", nil)
	req.Header.Set("Authorization", "Bearer secret")
	req.Header.Set("X-Forwarded-For", "203.0.113.10, 10.0.0.10")

	rec := httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if strings.Contains(body, "Bearer secret") {
		t.Fatal("expected authorization header to be redacted")
	}
	if !strings.Contains(body, "[redacted]") {
		t.Fatal("expected redacted marker in response")
	}
	if !strings.Contains(body, "203.0.113.10") {
		t.Fatal("expected forwarded address in response")
	}
}

type testWriter struct {
	t *testing.T
}

func (w testWriter) Write(p []byte) (int, error) {
	w.t.Log(strings.TrimSpace(string(p)))
	return len(p), nil
}
