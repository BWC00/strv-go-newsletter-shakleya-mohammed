package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/api/requestlog"
	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/api/middleware"
	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/util/logger"
)

func RegisterHTTPEndPoints(router *chi.Mux, l *logger.Logger, v *validator.Validate, postgresDB *gorm.DB) {
	newsletterAPI := New(l, v, postgresDB)
	middlewareHandlers := middleware.New(l, v)

	router.Route("/newsletters", func(r chi.Router) {
		// Protected routes
		r.Use(middlewareHandlers.Authentication)

		r.Method("GET", "/", requestlog.NewHandler(newsletterAPI.List, l))
		r.Method("GET", "/{id}", requestlog.NewHandler(newsletterAPI.Read, l))
		r.Method("DELETE", "/{id}", requestlog.NewHandler(newsletterAPI.Delete, l))

		r.Group(func(r chi.Router) {
			// Routes require validation
			r.Use(middlewareHandlers.ValidateNewsletter)

			r.Method("POST", "/", requestlog.NewHandler(newsletterAPI.Create, l))
			r.Method("PUT", "/{id}", requestlog.NewHandler(newsletterAPI.Update, l))
		})
	})
}
