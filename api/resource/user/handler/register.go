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
	userAPI := New(l, v, postgresDB)
	middlewareHandlers := middleware.New(l, v)

	router.Route("/users", func(r chi.Router) {
		r.Use(middlewareHandlers.ContentTypeJson)
		r.Use(middlewareHandlers.ValidateUser)

		r.Method("POST", "/login", requestlog.NewHandler(userAPI.Login, l))
		r.Method("POST", "/register", requestlog.NewHandler(userAPI.Register, l))
	})
}
