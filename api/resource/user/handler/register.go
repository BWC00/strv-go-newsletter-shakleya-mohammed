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

// RegisterHTTPEndPoints registers the HTTP endpoints for user registration and login.
func RegisterHTTPEndPoints(l *logger.Logger, v *validator.Validate, postgresDB *gorm.DB) http.Handler {
	// Create a new router instance
	r := chi.NewRouter()

	// Initialize the user API handler with the provided logger, validator, and database
	userAPI := New(l, v, postgresDB)

	// Initialize the middleware handlers with the provided logger and validator
	middlewareHandlers := middleware.New(l, v)

	// Apply the ValidateUser middleware to all routes registered on the router
	r.Use(middlewareHandlers.ValidateUser)

	// Handle the POST /login route
	// This route is responsible for logging in a user and returning a token
	r.Method("POST", "/login", requestlog.NewHandler(userAPI.Login, l))

	// Handle the POST /register route
	// This route is responsible for registering a new user and returning a token (automatically logged in)
	r.Method("POST", "/register", requestlog.NewHandler(userAPI.Register, l))

	// Return the configured router as an http.Handler
	return r
}
