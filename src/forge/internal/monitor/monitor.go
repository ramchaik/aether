package monitor

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Define a custom metric
var ProcessedMessages = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "forge_processed_messages_total",
		Help: "Total number of processed messages.",
	},
	[]string{"status"},
)

func init() {
	// Register the custom metric
	prometheus.MustRegister(ProcessedMessages)
}

func StartMetricsServer() {
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
