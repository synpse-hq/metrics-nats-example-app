package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/synpse-hq/metrics-nats-example-app/pkg/agent"
	"github.com/synpse-hq/metrics-nats-example-app/pkg/api"
	"github.com/synpse-hq/metrics-nats-example-app/pkg/logger"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()

	uri := "localhost:8080"
	val := os.Getenv("SYNPSE_APP_HOST")
	if val != "" {
		uri = val
	}

	cfg := api.Config{
		Host:                   uri,
		MetricsRefreshInterval: time.Second,
		NatsHost:               "nats://nats:4222",
		MetricsTopic:           "metrics",
	}

	logger := logger.GetLoggerInstance("", logger.ParseLogLevel(os.Getenv("SYNPSE_APP_LOG_LEVEL")))

	agent, err := agent.New(logger, cfg)
	if err != nil {
		return err
	}

	return agent.Run(ctx)

}
