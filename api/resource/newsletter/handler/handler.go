package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	vd "github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/api/resource/newsletter/repository"
	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/api/resource/newsletter"
	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/util/validator"
	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/util/logger"
	e "github.com/bwc00/strv-go-newsletter-shakleya-mohammed/util/err"
)

// API represents the handler for the newsletter API endpoints.
type API struct {
	logger     *logger.Logger
	validator  *vd.Validate
	repository *repository.Repository
}

// New creates a new instance of the API handler.
func New(logger *logger.Logger, validator *vd.Validate, postgresDB *gorm.DB) *API {
	return &API{
		logger:     logger,
		validator:  validator,
		repository: repository.NewRepository(postgresDB),
	}
}

// List godoc
//
//	@summary		List newsletters
//	@description	Handler function for listing newsletters
//	@tags			newsletters
//	@accept			json
//	@produce		json
//	@success		200		{array}		Newsletter
//	@failure		500		{object}	err.Error
//	@router			/newsletters [get]
func (a *API) List(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user ID from the request context
	userId := r.Context().Value(validator.UserKeyID).(uint32)

	// Retrieve the list of newsletters from the repository
	newsletters, err := a.repository.ListNewsletters(userId)
	if err != nil {
		a.logger.Error().Err(err).Msg("unable to list newsletters")
		e.ServerError(w, e.DataAccessFailure)
		return
	}

	// Check if there are no newsletters
	if newsletters == nil {
		fmt.Fprint(w, "[]")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Encode the newsletters into the response
	if err := json.NewEncoder(w).Encode(newsletters); err != nil {
		a.logger.Error().Err(err).Msg("unable to encode newsletters into response")
		e.ServerError(w, e.JsonEncodingFailure)
		return
	}

	// Return 200 OK status
	w.WriteHeader(http.StatusOK)
}

// Create godoc
//
//	@summary		Create newsletter
//	@description	Handler function for creating a newsletter
//	@tags			newsletters
//	@accept			json
//	@produce		json
//	@param			body	body		Newsletter	true	"Newsletter contents"
//	@success		201
//	@failure		400		{object}	err.Error
//	@failure		422		{object}	err.Errors
//	@failure		500		{object}	err.Error
//	@router			/newsletters [post]
func (a *API) Create(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user ID from the request context
	userId := r.Context().Value(validator.UserKeyID).(uint32)

	// Retrieve the newsletter from the request context
	newsletter := r.Context().Value(validator.ResourceKeyID).(*newsletter.Newsletter)
	newsletter.EditorId = userId

	// Create the newsletter using the repository
	result, err := a.repository.CreateNewsletter(newsletter)
	if err != nil {
		a.logger.Error().Err(err).Msg("unable to create newsletter")
		e.ServerError(w, e.DataCreationFailure)
		return
	}

	// Log successful creation of the newsletter
	a.logger.Info().Str("id", fmt.Sprintf("%d", result.ID)).Msg("newsletter created")

	// Return 200 OK status
	w.WriteHeader(http.StatusCreated)
}

// Read godoc
//
//	@summary		Read newsletter
//	@description	Handler function for reading a newsletter
//	@tags			newsletters
//	@accept			json
//	@produce		json
//	@param			id		path		string		true	"Newsletter ID"
//	@success		200		{object}	Newsletter
//	@failure		400		{object}	err.Error
//	@failure		404		{object}	err.Error
//	@failure		500		{object}	err.Error
//	@router			/newsletters/{id} [get]
func (a *API) Read(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user ID from the request context
	userId := r.Context().Value(validator.UserKeyID).(uint32)

	// Retrieve the newsletter ID from the URL parameter
	newsletterId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		a.logger.Error().Err(err).Msg("invalid newsletter id in url param")
		e.BadRequest(w, e.InvalidIdInUrlParam)
		return
	}

	// Read the newsletter using the repository
	newsletter, err2 := a.repository.ReadNewsletter(uint32(newsletterId), userId)
	if err2 != nil {
		if err2 == gorm.ErrRecordNotFound {
			// If the newsletter is not found, return a "Not Found" response
			a.logger.Error().Err(err2).Msg("newsletter not found")
			e.NotFoundError(w, e.ResourceNotFound)
			return
		}

		// If any other error occurs during the read, return a "Server Error" response
		a.logger.Error().Err(err2).Msg("unable to access newsletter")
		e.ServerError(w, e.DataAccessFailure)
		return
	}

	// Encode the newsletter into the response
	if err := json.NewEncoder(w).Encode(newsletter); err != nil {
		a.logger.Error().Err(err).Msg("unable to encode newsletter into response")
		e.ServerError(w, e.JsonEncodingFailure)
		return
	}

	// Return 200 OK status
	w.WriteHeader(http.StatusOK)
}

// Update godoc
//
//	@summary		Update newsletter
//	@description	Handler function for updating a newsletter
//	@tags			newsletters
//	@accept			json
//	@produce		json
//	@param			id		path		string		true	"Newsletter ID"
//	@param			body	body		Newsletter	true	"Updated newsletter contents"
//	@success		200
//	@failure		400		{object}	err.Error
//	@failure		404		{object}	err.Error
//	@failure		422		{object}	err.Errors
//	@failure		500		{object}	err.Error
//	@router			/newsletters/{id} [put]
func (a *API) Update(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user ID from the request context
	userId := r.Context().Value(validator.UserKeyID).(uint32)

	// Retrieve the newsletter from the request context
	newsletter := r.Context().Value(validator.ResourceKeyID).(*newsletter.Newsletter)

	// Retrieve the newsletter ID from the URL parameter
	newsletterId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		a.logger.Error().Err(err).Msg("invalid newsletter id in url param")
		e.BadRequest(w, e.InvalidIdInUrlParam)
		return
	}

	// Set the editor ID and newsletter ID
	newsletter.EditorId = userId
	newsletter.ID = uint32(newsletterId)

	// Update the newsletter using the repository
	if err := a.repository.UpdateNewsletter(newsletter); err != nil {
		// If the newsletter is not found, return a "Not Found" response
		if err == gorm.ErrRecordNotFound {
			a.logger.Error().Err(err).Msg("newsletter not found")
			e.NotFoundError(w, e.ResourceNotFound)
			return
		}

		// If any other error occurs during the update, return a "Server Error" response
		a.logger.Error().Err(err).Msg("unable to update newsletter")
		e.ServerError(w, e.DataUpdateFailure)
		return
	}

	// Log successful updating of the newsletter
	a.logger.Info().Str("id", fmt.Sprintf("%d", newsletter.ID)).Msg("newsletter updated")

	// Return 200 OK status
	w.WriteHeader(http.StatusOK)
}

// Delete godoc
//
//	@summary		Delete newsletter
//	@description	Handler function for deleting a newsletter
//	@tags			newsletters
//	@accept			json
//	@produce		json
//	@param			id		path		string	true	"Newsletter ID"
//	@success		204
//	@failure		400		{object}	err.Error
//	@failure		500		{object}	err.Error
//	@router			/newsletters/{id} [delete]
func (a *API) Delete(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user ID from the request context
	userId := r.Context().Value(validator.UserKeyID).(uint32)

	// Retrieve the newsletter ID from the URL parameter
	newsletterId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		a.logger.Error().Err(err).Msg("invalid newsletter id in url param")
		e.BadRequest(w, e.InvalidIdInUrlParam)
		return
	}

	// Delete the newsletter using the repository
	if err := a.repository.DeleteNewsletter(uint32(newsletterId), userId); err != nil {
		// If the newsletter is not found, return a "Not Found" response
		if err == gorm.ErrRecordNotFound {
			a.logger.Error().Err(err).Msg("newsletter not found")
			e.NotFoundError(w, e.ResourceNotFound)
			return
		}

		// If any other error occurs during deletion, return a "Server Error" response
		a.logger.Error().Err(err).Msg("unable to delete newsletter")
		e.ServerError(w, e.DataDeletionFailure)
		return
	}

	// Log successful deletion of the newsletter
	a.logger.Info().Str("id", fmt.Sprintf("%d", newsletterId)).Msg("newsletter deleted")

	// Return a "No Content" response indicating successful deletion
	w.WriteHeader(http.StatusNoContent)
}
