package server

import (
	"net/http"
	"fmt"

	"github.com/go-chi/chi/v5"

	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/api/middleware"
	e "github.com/bwc00/strv-go-newsletter-shakleya-mohammed/util/err"
	subscriptionAPI "github.com/bwc00/strv-go-newsletter-shakleya-mohammed/api/resource/subscription/handler"
	newsletterAPI "github.com/bwc00/strv-go-newsletter-shakleya-mohammed/api/resource/newsletter/handler"
	userAPI "github.com/bwc00/strv-go-newsletter-shakleya-mohammed/api/resource/user/handler"
)

// ROUTER AND MIDDLEWARE

// newRouter initializes a new chi Router for routing HTTP requests.
func (s *Server) newRouter() {
	s.router = chi.NewRouter()
}

// newMiddleware initializes a new middleware handler.
// It creates a new instance of the middleware handler using the server's logger and validator.
func (s *Server) newMiddleware() {
	s.middlewareHandler = middleware.New(s.logger, s.validator)
}

// setGlobalMiddleware adds global middleware to the router.
// It sets the content type to JSON for all requests handled by the router.
func (s *Server) setGlobalMiddleware() {
	s.router.Use(s.middlewareHandler.ContentTypeJson)
}

// registerHTTPEndPoints registers the HTTP endpoints for the server.
// It defines the routing logic for the server's API endpoints.
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