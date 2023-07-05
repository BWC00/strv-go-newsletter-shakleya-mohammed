package server

import (
	// Built-in
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	// Third party
	"gorm.io/gorm"
	"github.com/go-chi/chi/v5"
	vd "github.com/go-playground/validator/v10"
	"github.com/sendgrid/sendgrid-go"
	"firebase.google.com/go/v4/db"

	// User defined
	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/config"
	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/util/logger"
	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/api/middleware"
)

// SERVER ARCHITECTURE

// Server is the main struct representing the server and its components.
type Server struct {
	cfg               *config.Config
	postgresDB        *gorm.DB
	firebaseDB        *db.Ref
	sendGridClient    *sendgrid.Client
 	validator         *vd.Validate
	logger            *logger.Logger
	router            *chi.Mux
	middlewareHandler *middleware.Middleware
	httpServer        *http.Server
}

// New creates a new instance of the Server struct.
func New() *Server {
	return &Server{
		cfg:    config.New(),
	}
}


// INIT SERVER

// Initializes the server by setting up the logger, validator, database connections,
// SendGrid client, router, middleware, and HTTP endpoints registrations.
func (s *Server) Init() {
	s.newLogger()
	s.newValidator()

	s.newPostgresDB()
	s.newFirebaseDB()

	s.newSendGridClient()

	s.newRouter()
	s.newMiddleware()
	s.setGlobalMiddleware()
	s.registerHTTPEndPoints()
}


//RUN SERVER

// Run starts the server and listens for incoming requests.
func (s *Server) Run() {
	// Create an HTTP server instance with the server configuration settings
	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.cfg.Server.Port),
		Handler:      s.router,
		ReadTimeout:  s.cfg.Server.TimeoutRead,
		WriteTimeout: s.cfg.Server.TimeoutWrite,
		IdleTimeout:  s.cfg.Server.TimeoutIdle,
	}

	// Start the server in a separate goroutine
	go func() {
		start(s)
	}()

	// Perform graceful shutdown of the server
	if err := gracefulShutdown(s); err != nil {
		s.logger.Fatal().Err(err).Msg("Server shutdown failure")
	}
}

// start is an internal helper function that starts the HTTP server.
func start(s *Server) {
	s.logger.Info().Msgf("Starting server %v", s.httpServer.Addr)

	// Start the HTTP server and listen for incoming requests
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Fatal().Err(err).Msg("Server startup failure")
	}
}


//SHUTDOWN SERVER

// gracefulShutdown performs a graceful shutdown of the server.
// It handles the SIGINT and SIGTERM signals and shuts down the server gracefully.
func gracefulShutdown(s *Server) error {
	// Create a buffered channel to receive OS signals
	sigint := make(chan os.Signal, 1)

	// Notify the channel for specific signals
	signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Wait for a signal to be received on the channel
	<-sigint

	// Signal received, shutting down server
	s.logger.Info().Msgf("Shutting down server %v", s.httpServer.Addr)

	// Create a context with a timeout for the idle timeout duration
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.Server.TimeoutIdle)
	defer cancel()

	// Shutdown the HTTP server
	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Error().Err(err).Msg("Server shutdown failure")
	}

	// Close the PostgreSQL database connection
	sqlDB, err := s.postgresDB.DB()
	if err == nil {
		if err = sqlDB.Close(); err != nil {
			s.logger.Error().Err(err).Msg("postgres connection closing failure")
		}
	}

	return nil
}