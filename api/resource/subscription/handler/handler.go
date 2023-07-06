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

// API represents the handler for the subscription API endpoints.
type API struct {
	logger         *logger.Logger
	validator      *vd.Validate
	repository     *repository.Repository
	sendGridClient *sendgrid.Client
	cfg		   	   *config.EmailConfig
}

// New creates a new instance of the API handler.
func New(logger *logger.Logger, validator *vd.Validate, postgresDB *gorm.DB, firebaseDB *db.Ref, sendGridClient *sendgrid.Client, cfg *config.EmailConfig) *API {
	return &API{
		logger:     	logger,
		validator:  	validator,
		repository: 	repository.NewRepository(postgresDB, firebaseDB),
		sendGridClient: sendGridClient,
		cfg:			cfg,
	}
}

// List godoc
//
//	@summary		List subscriptions
//	@description	List subscriptions
//	@tags			subscriptions
//	@accept			json
//	@produce		json
//	@success		200
//	@failure		400		{object}	err.Error
//	@failure		422		{object}	err.Errors
//	@failure		500		{object}	err.Error
//	@router			/subscriptions [post]
func (a *API) List(w http.ResponseWriter, r *http.Request) {
	// Retrieve list of subscriptions
	subscriptions, err := a.repository.ListSubscriptions()
	if err != nil {
		a.logger.Error().Err(err).Msg("unable to list subscriptions")
		e.ServerError(w, e.DataAccessFailure)
		return
	}

	// Handle empty subscriptions case
	if subscriptions == nil {
		fmt.Fprint(w, "[]")

		// Return 200 OK status
		w.WriteHeader(http.StatusOK)
		return
	}

	// Encode subscriptions into response
	if err := json.NewEncoder(w).Encode(subscriptions); err != nil {
		a.logger.Error().Err(err).Msg("unable to encode subscription into response")
		e.ServerError(w, e.JsonEncodingFailure)
		return
	}

	// Return 200 OK status
	w.WriteHeader(http.StatusOK)
}

// Subscribe godoc
//
//	@summary		Create subscription
//	@description	Create subscription
//	@tags			subscriptions
//	@accept			json
//	@produce		json
//	@param			body	body		subscription.Subscription	true	"Subscription contents"
//	@success		201
//	@failure		400		{object}	err.Error
//	@failure		422		{object}	err.Errors
//	@failure		500		{object}	err.Error
//	@router			/subscriptions [post]
func (a *API) Subscribe(w http.ResponseWriter, r *http.Request) {
	// Extract subscription from the request context
	subscription := r.Context().Value(validator.ResourceKeyID).(*subscription.Subscription)

	// Create a new subscription
	subscriptionID, err := a.repository.Subscribe(subscription)
	if err != nil {
		// Check if the newsletter is not found
		if err == gorm.ErrRecordNotFound {
			a.logger.Error().Err(err).Msg("newsletter not found")
			e.NotFoundError(w, e.ResourceNotFound)
			return
		}

		// Handle other errors during subscription creation
		a.logger.Error().Err(err).Msg("Unable to create subscription")
		e.ServerError(w, e.DataCreationFailure)
		return
	}

	// Encode subscription ID into the response
	if err := json.NewEncoder(w).Encode(subscriptionID); err != nil {
		a.logger.Error().Err(err).Msg("unable to encode subscription id in response")
		e.ServerError(w, e.JsonEncodingFailure)
		return
	}

	// Get api version
	apiVersion := r.Context().Value(validator.ApiVersionKeyID).(string)

	// Prepare email subscription confirmation
	subject := "Subscribed to newsletter!"
	plainTextContent := "Subscribed to newsletter!"
	var unsubscribeURL string
	if r.TLS != nil {
		unsubscribeURL = fmt.Sprintf("https://%s/api/%s/subscriptions?id=%s", r.Host, apiVersion, subscriptionID)
	} else {
		unsubscribeURL = fmt.Sprintf("http://%s/api/%s/subscriptions?id=%s", r.Host, apiVersion, subscriptionID)
	}
	fmt.Println("URLLLLLLLLLLL:", unsubscribeURL)
	htmlContent := fmt.Sprintf("Subscribed! link to unsubscribe: <a href='%s'>unsubscribeYou</a>", unsubscribeURL)

	// Send email
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
		// Handle error if sending email fails
		a.logger.Error().Err(err).Msg("Unable to send email subscription confirmation")
		e.ServerError(w, e.SendingEmailFailure)
		return
	}

	// Log successful subscription
	a.logger.Info().Str("email", subscription.Email).Msg("subscribed to newsletter")

	// Return 201 Created status
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
//	@failure		500		{object}	err.Error
//	@router			/subscriptions [delete]
func (a *API) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	// Retrieve the subscription ID from the query parameter
	subscriptionID := r.URL.Query().Get("id")

	// Delete the subscription
	if err := a.repository.Unsubscribe(subscriptionID); err != nil {
		a.logger.Error().Err(err).Msg("unable to delete subscription")
		e.ServerError(w, e.DataDeletionFailure)
		return
	}

	// Log successful unsubscription
	a.logger.Info().Msg("unsubscribed to newsletter")

	// Write message confirmation
	fmt.Fprint(w, "Unsubscribed to newsletter!")

	// Return 204 No Content status
	w.WriteHeader(http.StatusNoContent)
}