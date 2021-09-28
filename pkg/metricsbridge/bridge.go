package metricsbridge

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"

	"github.com/synpse-hq/metrics-nats-example-app/pkg/api"
)

//go:embed config/*.yaml
var assets embed.FS

//metricsConfig stores the configuration of the metricsbridge application
type metricsConfig struct {
	Queries []struct {
		Query string `json:"query,omitempty"`
	} `json:"queries,omitempty"`

	DeviceID string `json:"deviceID,omitempty"`
}

func load() (*metricsConfig, error) {
	b, err := assets.ReadFile("config/config.yaml")
	if err != nil {
		return nil, err
	}
	var c metricsConfig
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *metricsConfig) defaultAndValidate() (err error) {
	if len(c.Queries) == 0 {
		return fmt.Errorf("must configure at least one query")
	}

	return
}

type MetricsBridge struct {
	log           *zap.Logger
	metricsConfig *metricsConfig
	nats          *nats.Conn

	config api.Config
}

func New(log *zap.Logger, cfg api.Config, nc *nats.Conn) (*MetricsBridge, error) {
	c, err := load()
	if err != nil {
		return nil, err
	}

	err = c.defaultAndValidate()
	if err != nil {
		return nil, err
	}

	b := &MetricsBridge{
		log:           log,
		metricsConfig: c,
		nats:          nc,
		config:        cfg,
	}

	return b, nil
}

func (m *MetricsBridge) Run(ctx context.Context) error {
	ticker := time.NewTicker(m.config.MetricsRefreshInterval)
	defer ticker.Stop()

	return m.runScraper(ctx)
}

func (m *MetricsBridge) runScraper(ctx context.Context) error {
	ticker := time.NewTicker(m.config.MetricsRefreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		}

		// we intentionally run after so metrics are populated. Small race condition with metrics package
		if err := m.runOnce(context.Background()); err != nil {
			m.log.Error("metrics bridge failed", zap.Error(err))
		}
	}
}

func (m *MetricsBridge) runOnce(ctx context.Context) error {
	register, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		return err
	}

	result := api.PrometheusMetric{
		Timestamp: time.Now().Unix(),
		Metrics:   make(map[string]float64, len(m.metricsConfig.Queries)),
	}

	// create list of queries we are interested to send back
	metrics := make(map[string]*dto.MetricFamily, len(m.metricsConfig.Queries))
	for _, query := range m.metricsConfig.Queries {
		metrics[query.Query] = nil
	}

	for _, r := range register {
		if r.Name != nil && *r.Name != "" {
			if _, ok := metrics[*r.Name]; ok {
				metrics[*r.Name] = r
			}
		}
	}

	for _, metric := range metrics {
		if metric.GetName() != "" {
			result.Metrics[metric.GetName()] = metric.Metric[0].Gauge.GetValue()
		}
	}

	data, err := json.Marshal(result)
	if err != nil {
		return err
	}

	// send to pub/sub locally for cloud SDKS to consume
	m.log.Info("publishing metrics")
	err = m.nats.Publish(m.config.MetricsTopic, data)
	if err != nil {
		return err
	}

	return nil
}
