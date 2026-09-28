package httpapi

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/robmoraes/dslab/assets/video-05/app/internal/config"
	"github.com/robmoraes/dslab/assets/video-05/app/internal/metadata"
)

type contextKey string

const requestIDKey contextKey = "request_id"

type Server struct {
	cfg       config.Config
	collector *metadata.Collector
	logger    *log.Logger
	sequence  atomic.Uint64
}

func NewServer(cfg config.Config, collector *metadata.Collector, logger *log.Logger) *Server {
	return &Server{
		cfg:       cfg,
		collector: collector,
		logger:    logger,
	}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleSnapshot)
	mux.HandleFunc("/metadata", s.handleSnapshot)
	mux.HandleFunc("/headers", s.handleHeaders)
	mux.HandleFunc("/healthz", s.handleHealth)
	mux.HandleFunc("/readyz", s.handleHealth)
	mux.HandleFunc("/version", s.handleVersion)
	mux.HandleFunc("/slow", s.handleSlow)

	return s.middleware(mux)
}

func (s *Server) handleSnapshot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/metadata" {
		writeError(w, r, http.StatusNotFound, "not found")
		return
	}
	if !isReadMethod(r.Method) {
		writeError(w, r, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	request := s.requestInfo(r)
	writeJSON(w, r, http.StatusOK, s.collector.Snapshot(r.Context(), request))
}

func (s *Server) handleHeaders(w http.ResponseWriter, r *http.Request) {
	if !isReadMethod(r.Method) {
		writeError(w, r, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	writeJSON(w, r, http.StatusOK, s.requestInfo(r))
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if !isReadMethod(r.Method) {
		writeError(w, r, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	writeJSON(w, r, http.StatusOK, map[string]string{
		"status": "ok",
		"app":    s.cfg.AppName,
	})
}

func (s *Server) handleVersion(w http.ResponseWriter, r *http.Request) {
	if !isReadMethod(r.Method) {
		writeError(w, r, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	writeJSON(w, r, http.StatusOK, map[string]any{
		"app":   s.cfg.AppName,
		"build": s.cfg.Build,
	})
}

func (s *Server) handleSlow(w http.ResponseWriter, r *http.Request) {
	if !isReadMethod(r.Method) {
		writeError(w, r, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	delay := delayFromQuery(r, 1000*time.Millisecond, 30*time.Second)
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-timer.C:
	case <-r.Context().Done():
		return
	}

	request := s.requestInfo(r)
	writeJSON(w, r, http.StatusOK, map[string]any{
		"delay_ms": delay.Milliseconds(),
		"snapshot": s.collector.Snapshot(r.Context(), request),
	})
}

func (s *Server) requestInfo(r *http.Request) metadata.RequestInfo {
	id, _ := r.Context().Value(requestIDKey).(string)
	remoteIP, remotePort := splitHostPort(r.RemoteAddr)
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	localAddr := ""
	if value := r.Context().Value(http.LocalAddrContextKey); value != nil {
		if addr, ok := value.(net.Addr); ok {
			localAddr = addr.String()
		}
	}

	return metadata.RequestInfo{
		ID:            id,
		Sequence:      s.sequence.Add(1),
		ReceivedAt:    time.Now().UTC(),
		Method:        r.Method,
		Path:          r.URL.Path,
		Query:         r.URL.RawQuery,
		Host:          r.Host,
		RemoteAddress: r.RemoteAddr,
		RemoteIP:      remoteIP,
		RemotePort:    remotePort,
		LocalAddress:  localAddr,
		Protocol:      r.Proto,
		Scheme:        scheme,
		UserAgent:     r.UserAgent(),
		Headers:       sanitizedHeaders(r.Header),
		LoadBalancer: metadata.LoadBalancerHeaders{
			Forwarded:       r.Header.Get("Forwarded"),
			ForwardedFor:    splitCSV(r.Header.Get("X-Forwarded-For")),
			ForwardedHost:   r.Header.Get("X-Forwarded-Host"),
			ForwardedProto:  r.Header.Get("X-Forwarded-Proto"),
			ForwardedPort:   r.Header.Get("X-Forwarded-Port"),
			RealIP:          r.Header.Get("X-Real-IP"),
			AmznTraceID:     r.Header.Get("X-Amzn-Trace-Id"),
			Via:             r.Header.Get("Via"),
			RequestID:       firstHeader(r.Header, "X-Request-Id", "X-Correlation-Id"),
			OriginalURI:     r.Header.Get("X-Original-URI"),
			OriginalForward: r.Header.Get("X-Original-Forwarded-For"),
		},
	}
}

func (s *Server) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestID := requestID(r)
		ctx := context.WithValue(r.Context(), requestIDKey, requestID)

		w.Header().Set("X-Request-Id", requestID)
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(lrw, r.WithContext(ctx))

		s.logger.Printf("request_id=%s method=%s path=%q status=%d duration=%s remote=%s", requestID, r.Method, r.URL.RequestURI(), lrw.statusCode, time.Since(start).Round(time.Millisecond), r.RemoteAddr)
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *loggingResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func writeJSON(w http.ResponseWriter, r *http.Request, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if r.Method == http.MethodHead {
		return
	}

	encoder := json.NewEncoder(w)
	if r.URL.Query().Get("pretty") == "1" || r.URL.Query().Get("pretty") == "true" {
		encoder.SetIndent("", "  ")
	}
	_ = encoder.Encode(payload)
}

func writeError(w http.ResponseWriter, r *http.Request, status int, message string) {
	writeJSON(w, r, status, map[string]any{
		"error":  message,
		"status": status,
	})
}

func isReadMethod(method string) bool {
	return method == http.MethodGet || method == http.MethodHead
}

func requestID(r *http.Request) string {
	if value := strings.TrimSpace(firstHeader(r.Header, "X-Request-Id", "X-Correlation-Id")); value != "" {
		return value
	}

	var bytes [8]byte
	if _, err := rand.Read(bytes[:]); err == nil {
		return hex.EncodeToString(bytes[:])
	}

	return strconv.FormatInt(time.Now().UnixNano(), 36)
}

func sanitizedHeaders(headers http.Header) map[string]string {
	result := make(map[string]string, len(headers))
	for key, values := range headers {
		canonical := http.CanonicalHeaderKey(key)
		if isSensitiveHeader(canonical) {
			result[canonical] = "[redacted]"
			continue
		}
		result[canonical] = strings.Join(values, ", ")
	}
	return result
}

func isSensitiveHeader(key string) bool {
	switch strings.ToLower(key) {
	case "authorization", "cookie", "set-cookie", "x-api-key", "x-auth-token":
		return true
	default:
		return false
	}
}

func firstHeader(headers http.Header, keys ...string) string {
	for _, key := range keys {
		if value := strings.TrimSpace(headers.Get(key)); value != "" {
			return value
		}
	}
	return ""
}

func splitCSV(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}

	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}

func splitHostPort(addr string) (string, string) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return addr, ""
	}
	return host, port
}

func delayFromQuery(r *http.Request, fallback time.Duration, max time.Duration) time.Duration {
	raw := strings.TrimSpace(r.URL.Query().Get("ms"))
	if raw == "" {
		return fallback
	}

	ms, err := strconv.Atoi(raw)
	if err != nil || ms < 0 {
		return fallback
	}

	delay := time.Duration(ms) * time.Millisecond
	if delay > max {
		return max
	}
	return delay
}
