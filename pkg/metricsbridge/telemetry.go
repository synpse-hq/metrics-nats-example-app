package metricsbridge

import (
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	prometheus.MustRegister(metricsBridgeErrorsCounter)
	prometheus.MustRegister(metricsBridgeProcessingDurationSummary)
	prometheus.MustRegister(metricsBridgeMetricsTransferredCounter)

}

var (
	metricsBridgeErrorsCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "metrics_bridge_errors_total",
			Help: "Total number of errors.",
		},
	)

	metricsBridgeProcessingDurationSummary = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "metrics_bridge_processing_duration_seconds",
			Help: "The duration of metrics bridge query processing.",
		},
	)

	metricsBridgeMetricsTransferredCounter = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "metrics_bridge_metrics_transferred_total",
			Help: "Total number of metrics transferred.",
		},
	)
)
