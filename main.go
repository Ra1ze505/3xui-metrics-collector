package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Configuration loaded: Host=%s, Port=%s, BasePath=%s, UseTLS=%t",
		config.XUIHost, config.XUIPort, config.XUIBasePath, config.XUIUseTLS)

	client := NewXUIClient(config)

	log.Printf("Attempting to login to X-UI panel...")
	if err := client.Login(); err != nil {
		log.Fatalf("Failed to login to X-UI: %v", err)
	}
	log.Printf("Successfully logged in to X-UI panel")

	collector := NewMetricsCollector(client)
	prometheus.MustRegister(collector)

	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Starting metrics server on :2112")
	if err := http.ListenAndServe(":2112", nil); err != nil {
		log.Fatalf("Failed to start metrics server: %v", err)
	}
}
