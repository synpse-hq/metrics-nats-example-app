package agent

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/synpse-hq/metrics-nats-example-app/pkg/api"
	"github.com/synpse-hq/metrics-nats-example-app/pkg/metrics"
	"github.com/synpse-hq/metrics-nats-example-app/pkg/metricsbridge"
	"github.com/synpse-hq/metrics-nats-example-app/pkg/service"
)

type Agent struct {
	service       service.Interface
	metrics       metrics.Metrics
	metricsbridge *metricsbridge.MetricsBridge
	nats          *nats.Conn
}

func New(log *zap.Logger, cfg api.Config) (*Agent, error) {

	// Connect to a server

	var connected bool
	var nc *nats.Conn
	for !connected {
		var err error
		nc, err = nats.Connect(cfg.NatsHost)
		if err != nil {
			log.Warn("failed to connect to nats", zap.Error(err))
		} else {
			log.Info("connect to nats!")
			connected = true
		}
		time.Sleep(time.Second)

	}

	mb, err := metricsbridge.New(log, cfg, nc)
	if err != nil {
		return nil, err
	}

	return &Agent{
		service:       service.New(log, cfg),
		metrics:       metrics.New(log, cfg),
		metricsbridge: mb,
		nats:          nc,
	}, nil
}

func (a Agent) Run(ctx context.Context) error {
	g := &errgroup.Group{}

	g.Go(func() error {
		return a.service.Run(ctx)
	})

	g.Go(func() error {
		return a.metrics.Run(ctx)
	})

	g.Go(func() error {
		return a.metricsbridge.Run(ctx)
	})

	return g.Wait()
}
