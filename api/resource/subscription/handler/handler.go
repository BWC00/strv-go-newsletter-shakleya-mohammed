package handler

import (
	"encoding/json"
	"net/http"
	"fmt"

	vd "github.com/go-playground/validator/v10"
	"github.com/sendgrid/sendgrid-go"
	"firebase.google.com/go/v4/db"
	"gorm.io/gorm"

	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/api/resource/subscription/repository"
	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/api/resource/subscription"
	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/util/validator"
	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/util/logger"
	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/util/email"
	e "github.com/bwc00/strv-go-newsletter-shakleya-mohammed/util/err"
	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/config"
)

type API struct {
	logger         *logger.Logger
	validator      *vd.Validate
	repository     *repository.Repository
	sendGridClient *sendgrid.Client
	cfg		   	   *config.EmailConfig
}

func New(logger *logger.Logger, validator *vd.Validate, postgresDB *gorm.DB, firebaseDB *db.Ref, sendGridClient *sendgrid.Client, cfg *config.EmailConfig) *API {
	return &API{
		logger:     	logger,
		validator:  	validator,
		repository: 	repository.NewRepository(postgresDB, firebaseDB),
		sendGridClient: sendGridClient,
		cfg:			cfg,
	}
}

// Subscribe godoc
//
//	@summary		Create subscription
//	@description	Create subscription
//	@tags			subscriptions
//	@accept			json
//	@produce		json
//	@param			body	body		Subscription	true	"Subscription contents"
//	@success		201
//	@failure		400		{object}	err.Error
//	@failure		422		{object}	err.Errors
//	@failure		500		{object}	err.Error
//	@router			/subscriptions [post]
func (a *API) Subscribe(w http.ResponseWriter, r *http.Request) {
	
	subscription := r.Context().Value(validator.KeyID).(*subscription.Subscription)

	subscriptionID, err := a.repository.Subscribe(subscription)
	if err != nil {
		a.logger.Error().Err(err).Msg("Unable to create subscription")
		e.ServerError(w, e.DataCreationFailure)
		return
	}

	if err := json.NewEncoder(w).Encode(subscriptionID); err != nil {
		a.logger.Error().Err(err).Msg("unable to encode subscription id in response")
		e.ServerError(w, e.JsonEncodingFailure)
		return
	}

	subject := "Subscribed to newsletter!"
	plainTextContent := "Subscribed to newsletter!"
	htmlContent := fmt.Sprintf("Subscribed! link to unsubscribe: <a href='http://localhost:8080/api/v1/unsubscribe?subscriptionid=%s'>unsubscribeYou</a>", subscriptionID)

	if err := email.Send(
		a.sendGridClient,
		a.cfg.SendGrid.SendFromName,
		a.cfg.SendGrid.SendFromAddress,
		subject,
		email.ExtractEmailUsername(subscription.Email),
		subscription.Email,
		plainTextContent,
		htmlContent,
	); err != nil {
		a.logger.Error().Err(err).Msg("Unable to send notification to email")
		e.ServerError(w, e.SendingEmailFailure)
	}

	a.logger.Info().Str("email", subscription.Email).Msg("subscribed to newsletter")
	w.WriteHeader(http.StatusCreated)
}

// Unsubscribe godoc
//
//	@summary		Delete subscription
//	@description	Delete subscription
//	@tags			subscriptions
//	@accept			json
//	@produce		json
//	@param			id		path		string	true	"Subscription ID"
//	@success		204
//	@failure		404		{object}	err.Error
//	@router			/subscriptions [delete]
func (a *API) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	subscriptionID := r.URL.Query().Get("id")
	if err := a.repository.Unsubscribe(subscriptionID); err != nil {
		a.logger.Error().Err(err).Msg("already unsubscribed")
		e.NotFoundErrors(w, e.ResourceNotFound)
		return
	}

	a.logger.Info().Msg("unsubscribed to newsletter")
}