package api

import "time"

type PrometheusMetric struct {
	ID        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
	ProjectID string `json:"projectId"`
	DeviceID  string `json:"deviceId"`

	Metrics map[string]float64 `json:"metrics"`
	Labels  map[string]string  `json:"labels"`
}

type Config struct {
	Host                   string        `json:"host"`
	MetricsRefreshInterval time.Duration `json:"metricsRefreshInterval"`
	NatsHost               string        `json:"natsHost"`
	MetricsTopic           string        `json:"metricsTopic"`
}
