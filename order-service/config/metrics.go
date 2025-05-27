// config/metrics.go
package config

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	RequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"service", "method", "status"},
	)
)

// InitMetrics запускает HTTP-сервер для экспорта метрик на указанном порту
func InitMetrics(port string, serviceName string) {
	prometheus.MustRegister(RequestCount)

	http.Handle("/metrics", promhttp.Handler())

	go func() {
		log.Printf("[%s] Prometheus metrics server starting on :%s/metrics", serviceName, port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatalf("Metrics server failed: %v", err)
		}
	}()
}
