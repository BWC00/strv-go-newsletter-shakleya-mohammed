package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"gorm.io/gorm"
	"github.com/go-chi/chi/v5"

	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/config"
	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/util/logger"
	databases "github.com/bwc00/strv-go-newsletter-shakleya-mohammed/database"

)

type Server struct {
	cfg            *config.Config
	postgresDB     *gorm.DB
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
	s.newPostgresDB()
}

func (s *Server) newLogger() {
	s.logger = logger.New(s.cfg.Server.Debug)
}

func (s *Server) newRouter() {
	s.router = chi.NewRouter()
}

func (s *Server) newPostgresDB() {
	var err error
	if s.postgresDB, err = databases.NewPostgresDB(&s.cfg.DB.RDBMS, s.logger); err != nil {
		s.logger.Fatal().Err(err).Msg("error initializing postgres database")
	}
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

	_ = gracefulShutdown(s)
}

func start(s *Server) {
	s.logger.Info().Msgf("Starting server %v", s.httpServer.Addr)
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Fatal().Err(err).Msg("Server startup failure")
	}
}

func gracefulShutdown(s *Server) error {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-sigint

	s.logger.Info().Msgf("Shutting down server %v", s.httpServer.Addr)

	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.Server.TimeoutIdle)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Error().Err(err).Msg("Server shutdown failure")
	}

	sqlDB, err := s.postgresDB.DB()
	if err == nil {
		if err = sqlDB.Close(); err != nil {
			s.logger.Error().Err(err).Msg("postgres connection closing failure")
		}
	}

	return nil
}