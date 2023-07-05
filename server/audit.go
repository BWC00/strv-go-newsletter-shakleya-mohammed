// LOGGING AND VALIDATION

func (s *Server) newLogger() {
	s.logger = logger.New(s.cfg.Server.Debug)
}

func (s *Server) newValidator() {
	s.validator = validator.New()
}
