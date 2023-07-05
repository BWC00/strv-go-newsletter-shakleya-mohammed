package server

import (
	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/util/logger"
	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/util/validator"
)

// LOGGING AND VALIDATION

// newLogger creates a new logger instance for the server.
// It initializes the logger with the debug mode specified in the server configuration.
func (s *Server) newLogger() {
	s.logger = logger.New(s.cfg.Server.Debug)
}

// newValidator creates a new validator instance for the server.
// It initializes the validator with default settings and returns it.
func (s *Server) newValidator() {
	s.validator = validator.New()
}
