package server

import (
	databases "github.com/bwc00/strv-go-newsletter-shakleya-mohammed/database"
)

// DATABASE SERVICES

// newPostgresDB initializes a new PostgreSQL database connection for the server.
// It creates a new PostgreSQL database connection based on the configuration settings.
// If any error occurs during the initialization, it logs a fatal error and terminates the server.
func (s *Server) newPostgresDB() {
	var err error
	if s.postgresDB, err = databases.NewPostgresDB(&s.cfg.DB.RDBMS, s.logger); err != nil {
		s.logger.Fatal().Err(err).Msg("error initializing postgres database")
	}
}

// newFirebaseDB initializes a new Firebase Realtime Database connection for the server.
// It creates a new Firebase Realtime Database connection based on the configuration settings.
// If any error occurs during the initialization, it logs a fatal error and terminates the server.
func (s *Server) newFirebaseDB() {
	var err error
	if s.firebaseDB, err = databases.NewFirebaseDB(&s.cfg.DB.Firebase); err != nil {
		s.logger.Fatal().Err(err).Msg("error initializing firebase database")
	}
}