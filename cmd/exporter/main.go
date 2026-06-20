package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/andrejmatveev/3xui-metrics-collector/internal/collector"
	"github.com/andrejmatveev/3xui-metrics-collector/internal/config"
	"github.com/andrejmatveev/3xui-metrics-collector/internal/poller"
	"github.com/andrejmatveev/3xui-metrics-collector/internal/xui"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	configPath := flag.String("config", envOrDefault("CONFIG_PATH", "config.yaml"), "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var pollers []*poller.Poller
	for _, panel := range cfg.Panels {
		client := xui.NewClient(panel.BaseURL, panel.APIToken, cfg.RequestTimeout, panel.InsecureSkipVerify)
		p := poller.New(panel.Name, client, cfg.PollInterval, cfg.RequestTimeout, panel.CollectOutbounds)
		pollers = append(pollers, p)
		go p.Start(ctx)
	}

	var sources []collector.SnapshotSource
	for _, p := range pollers {
		sources = append(sources, p)
	}
	col := collector.New(sources...)
	prometheus.MustRegister(col)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	srv := &http.Server{
		Addr:    cfg.ListenAddr,
		Handler: mux,
	}

	go func() {
		log.Printf("listening on %s", cfg.ListenAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server: %v", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("shutting down...")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
