package server

import (
	databases "github.com/bwc00/strv-go-newsletter-shakleya-mohammed/database"
)

// DATABASE SERVICES

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