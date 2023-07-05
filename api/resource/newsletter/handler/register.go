package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/api/requestlog"
	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/api/middleware"
	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/util/logger"
)

// RegisterHTTPEndPoints registers the HTTP endpoints for newsletters and returns an http.Handler.
// It configures the routes, applies middleware, and associates the handler functions with the respective HTTP methods and paths.
// The provided logger, validator, and postgresDB are used for request logging, input validation, and database operations respectively.
// The returned http.Handler can be used to serve the HTTP endpoints.
func RegisterHTTPEndPoints(l *logger.Logger, v *validator.Validate, postgresDB *gorm.DB) http.Handler {
	// Create a new router instance
	r := chi.NewRouter()

	// Initialize the user API handler with the provided logger, validator, and database
	newsletterAPI := New(l, v, postgresDB)

	// Initialize the middleware handlers with the provided logger and validator
	middlewareHandlers := middleware.New(l, v)

	// Apply the Authentication middleware to all routes
	r.Use(middlewareHandlers.Authentication)

	// Define routes for GET and DELETE requests
	r.Method("GET", "/", requestlog.NewHandler(newsletterAPI.List, l))
	r.Method("GET", "/{id}", requestlog.NewHandler(newsletterAPI.Read, l))
	r.Method("DELETE", "/{id}", requestlog.NewHandler(newsletterAPI.Delete, l))

	// Group routes that require validation middleware
	r.Group(func(r chi.Router) {
		r.Use(middlewareHandlers.ValidateNewsletter)

		// Define routes for POST and PUT requests
		r.Method("POST", "/", requestlog.NewHandler(newsletterAPI.Create, l))
		r.Method("PUT", "/{id}", requestlog.NewHandler(newsletterAPI.Update, l))
	})

	return r
}
