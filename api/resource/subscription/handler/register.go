package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/sendgrid/sendgrid-go"
	"firebase.google.com/go/v4/db"
	"gorm.io/gorm"

	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/api/requestlog"
	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/api/middleware"
	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/util/logger"
	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/config"
)

// RegisterHTTPEndPoints registers the HTTP endpoints for subscriptions and returns the HTTP handler.
// It configures routing and middleware for the endpoints, and associates the handler functions with the respective HTTP paths.
// The provided logger, validator, and postgresDB are used for request logging, input validation, and database operations respectively.
// The returned http.Handler can be used to serve the HTTP endpoints.
func RegisterHTTPEndPoints(l *logger.Logger, v *validator.Validate, postgresDB *gorm.DB, firebaseDB *db.Ref, sendGridClient *sendgrid.Client, cfg *config.EmailConfig) http.Handler {
	// Create a new router instance
	r := chi.NewRouter()

	// Create a new instance of API handler for subscriptions
	subscriptionAPI := New(l, v, postgresDB, firebaseDB, sendGridClient, cfg)

	// Create a new instance of the middleware handlers
	middlewareHandlers := middleware.New(l, v)

	// GET /all
	// Handle the list subscriptions endpoint
	// Log the request using requestlog middleware
	r.Method("GET", "/all", requestlog.NewHandler(subscriptionAPI.List, l))

	// GET /
	// Handle the unsubscribe endpoint
	// Log the request using requestlog middleware
	r.Method("GET", "/", requestlog.NewHandler(subscriptionAPI.Unsubscribe, l))

	// Group the routes that require subscription validation middleware
	r.Group(func(r chi.Router) {
		// Apply the subscription validation middleware to this route group
		r.Use(middlewareHandlers.ValidateSubscription)

		// POST /
		// Handle the subscribe endpoint
		// Log the request using requestlog middleware
		r.Method("POST", "/", requestlog.NewHandler(subscriptionAPI.Subscribe, l))
	})

	// Return the router as the HTTP handler
	return r
}
