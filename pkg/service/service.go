package service

import (
	"context"
	"net"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/synpse-hq/metrics-nats-example-app/pkg/api"
)

var _ Interface = &Service{}

type Interface interface {
	Run(context.Context) error
	Shutdown(context.Context) error
}

type Service struct {
	log    *zap.Logger
	server *http.Server
	router *mux.Router

	config api.Config
}

func New(log *zap.Logger, cfg api.Config) *Service {
	s := &Service{
		log: log,

		config: cfg,
	}

	s.router = s.setupRouter()

	s.router.Handle("/metrics", promhttp.Handler())

	s.router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	s.router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	s.router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	s.router.HandleFunc("/debug/pprof/trace", pprof.Trace)
	s.router.PathPrefix("/debug/pprof/").HandlerFunc(pprof.Index)

	s.server = &http.Server{
		Addr:    cfg.Host,
		Handler: s.router,
	}

	return s
}

func (s *Service) Run(ctx context.Context) error {
	s.log.Info("Starting Agent Local API", zap.String("url", s.config.Host))
	defer s.log.Info("Agent API Local shutdown completed", zap.String("url", s.config.Host))

	listener, err := net.Listen("tcp", s.config.Host)
	if err != nil {
		return err
	}
	defer listener.Close()

	return s.server.Serve(listener)
}

func (s *Service) Shutdown(ctx context.Context) error {
	s.log.Info("Agent API shutdown started")
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	err := s.server.Shutdown(ctx)
	if err != nil {
		return errors.Wrap(err, "shutdown remote server")
	}
	return nil
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Service) setupRouter() *mux.Router {
	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
	})

	return r
}
