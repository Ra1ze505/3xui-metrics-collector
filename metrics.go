package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type MetricsCollector struct {
	client *XUIClient
	up     *prometheus.GaugeVec
	down   *prometheus.GaugeVec
	mu     sync.Mutex
}

func NewMetricsCollector(client *XUIClient) *MetricsCollector {
	return &MetricsCollector{
		client: client,
		up: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "xui_client_up_bytes",
				Help: "Total uploaded bytes for each client",
				ConstLabels: prometheus.Labels{
					"type": "integer",
				},
			},
			[]string{"email", "client_id", "inbound_id", "inbound_remark"},
		),
		down: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "xui_client_down_bytes",
				Help: "Total downloaded bytes for each client",
				ConstLabels: prometheus.Labels{
					"type": "integer",
				},
			},
			[]string{"email", "client_id", "inbound_id", "inbound_remark"},
		),
	}
}

func (m *MetricsCollector) Describe(ch chan<- *prometheus.Desc) {
	m.up.Describe(ch)
	m.down.Describe(ch)
}

func (m *MetricsCollector) Collect(ch chan<- prometheus.Metric) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.up.Reset()
	m.down.Reset()

	log.Printf("Starting metrics collection...")
	inbounds, err := m.client.GetInbounds()
	if err != nil {
		log.Printf("Error getting inbounds: %v", err)
		return
	}

	clientCount := 0
	for _, inbound := range inbounds {
		for _, client := range inbound.ClientStats {
			if !client.Enable {
				continue
			}

			clientCount++
			labels := prometheus.Labels{
				"email":          client.Email,
				"client_id":      fmt.Sprintf("%d", client.ID),
				"inbound_id":     fmt.Sprintf("%d", inbound.ID),
				"inbound_remark": inbound.Remark,
			}

			upValue := float64(client.Up)
			downValue := float64(client.Down)

			m.up.With(labels).Set(upValue)
			m.down.With(labels).Set(downValue)
		}
	}

	log.Printf("Collected metrics for %d active clients", clientCount)
	m.up.Collect(ch)
	m.down.Collect(ch)
}
