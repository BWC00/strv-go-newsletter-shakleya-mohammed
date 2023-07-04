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
	vd "github.com/go-playground/validator/v10"
	"github.com/sendgrid/sendgrid-go"
	"firebase.google.com/go/v4/db"

	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/config"
	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/util/logger"
	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/util/validator"
	databases "github.com/bwc00/strv-go-newsletter-shakleya-mohammed/database"

)

type Server struct {
	cfg            *config.Config
	postgresDB     *gorm.DB
	firebaseDB     *db.Ref
	sendGridClient *sendgrid.Client
 	validator      *vd.Validate
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
	s.newValidator()
	s.newRouter()
	s.newPostgresDB()
	s.newFirebaseDB()
	s.newSendGridClient()
}

func (s *Server) newLogger() {
	s.logger = logger.New(s.cfg.Server.Debug)
}

func (s *Server) newValidator() {
	s.validator = validator.New()
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

func (s *Server) newFirebaseDB() {
	var err error
	if s.firebaseDB, err = databases.NewFirebaseDB(&s.cfg.DB.Firebase); err != nil {
		s.logger.Fatal().Err(err).Msg("error initializing firebase database")
	}
}

func (s *Server) newSendGridClient() {
	s.sendGridClient = sendgrid.NewSendClient(s.cfg.Email.SendGrid.ApiKey)
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