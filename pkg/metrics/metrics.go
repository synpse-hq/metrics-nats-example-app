package metrics

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/synpse-hq/metrics-nats-example-app/pkg/api"
)

type Metrics struct {
	log    *zap.Logger
	config api.Config

	cpuCache *cpuCache
}

type cpuCache struct {
	User    uint64
	System  uint64
	Idle    uint64
	Guest   uint64
	Nice    uint64
	Iowait  uint64
	Total   uint64
	Softirq uint64
	Steal   uint64

	// CPU count
	Count int
}

func New(log *zap.Logger, cfg api.Config) Metrics {
	return Metrics{
		log:    log,
		config: cfg,
	}
}

func (m *Metrics) Run(ctx context.Context) error {
	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	fMap := []func(context.Context){
		m.updateMemory,
		m.updateCPU,
	}

	// on the first run we will run function 2 times with interval of 10 seconds
	// to have cold start metric values. And we move to usuall ticker after this.
	for _, f := range fMap {
		f(ctx)
	}

	select {
	case <-ctx.Done():
		return nil
	case <-time.After(10 * time.Second):
		// continue
	}

	for _, f := range fMap {
		f(ctx)
	}

	for {
		select {
		case <-ticker.C:

			for _, f := range fMap {
				f(ctx)
			}

		case <-ctx.Done():
			return nil
		}
	}
}
