package handler

import (
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

func RegisterHTTPEndPoints(router *chi.Mux, l *logger.Logger, v *validator.Validate, postgresDB *gorm.DB, firebaseDB *db.Ref, sendGridClient *sendgrid.Client, cfg *config.EmailConfig) {
	subscriptionAPI := New(l, v, postgresDB, firebaseDB, sendGridClient, cfg)
	middlewareHandlers := middleware.New(l, v)

	router.Route("/subscriptions", func(r chi.Router) {
		r.Use(middlewareHandlers.ContentTypeJson)

		r.Method("GET", "/unsubscribe", requestlog.NewHandler(subscriptionAPI.Unsubscribe, l))

		r.Group(func(r chi.Router) {
			r.Use(middlewareHandlers.ValidateSubscription)

			r.Method("POST", "/subscribe", requestlog.NewHandler(subscriptionAPI.Subscribe, l))
		})
	})
}
