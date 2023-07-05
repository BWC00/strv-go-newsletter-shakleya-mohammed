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

func RegisterHTTPEndPoints(l *logger.Logger, v *validator.Validate, postgresDB *gorm.DB) http.Handler {
	r := chi.NewRouter()
	userAPI := New(l, v, postgresDB)
	middlewareHandlers := middleware.New(l, v)

	r.Use(middlewareHandlers.ValidateUser)

	r.Method("POST", "/login", requestlog.NewHandler(userAPI.Login, l))
	r.Method("POST", "/register", requestlog.NewHandler(userAPI.Register, l))

	return r
}
