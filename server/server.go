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

)

type Server struct {
	cfg            *config.Config
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
	s.newRouter()
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
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return
	}
}

func gracefulShutdown(ctx context.Context, s *Server) error {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-sigint

	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.Server.TimeoutIdle)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}