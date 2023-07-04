package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"

	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/config"
	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/util/logger"

)

type Server struct {
	cfg            *config.Config
	logger         *logger.Logger
	router         *chi.Mux
	httpServer     *http.Server
}

func New() *Server {
	return &Server{
		cfg:    config.New(),
	}
}


// INIT SERVER

func (s *Server) Init() {
	s.newLogger()
	s.newRouter()
}

func (s *Server) newLogger() {
	s.logger = logger.New(s.cfg.Server.Debug)
}

func (s *Server) newRouter() {
	s.router = chi.NewRouter()
}


//START SERVER

func (s *Server) Run() {
	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.cfg.Server.Port),
		Handler:      s.router,
		ReadTimeout:  s.cfg.Server.TimeoutRead,
		WriteTimeout: s.cfg.Server.TimeoutWrite,
		IdleTimeout:  s.cfg.Server.TimeoutIdle,
	}

	go func() {
		start(s)
	}()

	_ = gracefulShutdown(context.Background(), s)
}

func start(s *Server) {
	s.logger.Info().Msgf("Starting server %v", s.httpServer.Addr)
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Fatal().Err(err).Msg("Server startup failure")
	}
}

func gracefulShutdown(ctx context.Context, s *Server) error {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-sigint

	s.logger.Info().Msgf("Shutting down server %v", s.httpServer.Addr)

	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.Server.TimeoutIdle)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Error().Err(err).Msg("Server shutdown failure")
	}

	return nil
}