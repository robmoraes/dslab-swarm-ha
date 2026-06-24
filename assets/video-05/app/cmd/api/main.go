package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/robmoraes/dslab/assets/video-05/app/internal/config"
	"github.com/robmoraes/dslab/assets/video-05/app/internal/httpapi"
	"github.com/robmoraes/dslab/assets/video-05/app/internal/metadata"
)

var (
	version   = "dev"
	commit    = "local"
	buildDate = "unknown"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.LUTC)

	cfg := config.Load(config.BuildInfo{
		Version: version,
		Commit:  commit,
		Date:    buildDate,
	})

	collector := metadata.NewCollector(cfg)
	api := httpapi.NewServer(cfg, collector, logger)

	server := &http.Server{
		Addr:              cfg.ListenAddr,
		Handler:           api.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Printf("starting app=%s version=%s addr=%s env=%s", cfg.AppName, cfg.Build.Version, cfg.ListenAddr, cfg.Environment)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("http server failed: %v", err)
		}
	}()

	<-ctx.Done()
	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Print("shutting down")
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Fatalf("graceful shutdown failed: %v", err)
	}
}
