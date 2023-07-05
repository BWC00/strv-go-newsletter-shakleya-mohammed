// ROUTER AND MIDDLEWARE

func (s *Server) newRouter() {
	s.router = chi.NewRouter()
}

func (s *Server) newMiddleware() {
	s.middlewareHandler = middleware.New(s.logger, s.validator)
}

func (s *Server) setGlobalMiddleware() {
	s.router.Use(s.middlewareHandler.ContentTypeJson)
}

func (s *Server) registerHTTPEndPoints() {
	// Live route
	s.router.Get("/live", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("."))
	})

	// Not found route
	s.router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		e.NotFoundError(w, e.EndpointNotFound)
	})

	// API version 1
	s.router.Route(fmt.Sprintf("/api/v%d",s.cfg.Server.Major), func(r chi.Router) {
		r.Use(s.middlewareHandler.ApiVersionCtx(fmt.Sprintf("v%d",s.cfg.Server.Major)))
		r.Mount("/users", userAPI.RegisterHTTPEndPoints(s.logger, s.validator, s.postgresDB))
		r.Mount("/newsletters", newsletterAPI.RegisterHTTPEndPoints(s.logger, s.validator, s.postgresDB))
		r.Mount("/subscriptions", subscriptionAPI.RegisterHTTPEndPoints(s.logger, s.validator, s.postgresDB, s.firebaseDB, s.sendGridClient, &s.cfg.Email))
	})
}