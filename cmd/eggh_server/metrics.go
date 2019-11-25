package main

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	activeRooms = promauto.NewGauge(prometheus.GaugeOpts{
		Help:      "Active rooms",
		Namespace: "egghtogether",
		Subsystem: "roomserver",
		Name:      "active_rooms",
	})
	activeConnections = promauto.NewGauge(prometheus.GaugeOpts{
		Help:      "Active Connections",
		Namespace: "egghtogether",
		Subsystem: "roomserver",
		Name:      "active_connections",
	})
)

func recordMetrics() {
	go func() {
		for {
			activeRooms.Set(float64(manager.TotalRooms()))
			time.Sleep(1 * time.Second)
		}
	}()
	go func() {
		for {
			activeConnections.Set(float64(manager.TotalConnections()))
			time.Sleep(1 * time.Second)
		}
	}()
}
