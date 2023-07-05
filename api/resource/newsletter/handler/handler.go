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

type API struct {
	logger     *logger.Logger
	validator  *vd.Validate
	repository *repository.Repository
}

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
//	@description	List newsletters
//	@tags			newsletters
//	@accept			json
//	@produce		json
//	@success		200		{array}		Newsletter
//	@failure		500		{object}	err.Error
//	@router			/newsletters [get]
func (a *API) List(w http.ResponseWriter, r *http.Request) {

	userId := r.Context().Value(validator.UserKeyID).(uint32)

	newsletters, err := a.repository.ListNewsletters(userId)
	if err != nil {
		a.logger.Error().Err(err).Msg("unable to list newsletters")
		e.ServerError(w, e.DataAccessFailure)
		return
	}

	if newsletters == nil {
		fmt.Fprint(w, "[]")
		w.WriteHeader(http.StatusOK)
		return
	}

	if err := json.NewEncoder(w).Encode(newsletters); err != nil {
		a.logger.Error().Err(err).Msg("unable to encode newsletters into response")
		e.ServerError(w, e.JsonEncodingFailure)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Create godoc
//
//	@summary		Create newsletter
//	@description	Create newsletter
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

	userId := r.Context().Value(validator.UserKeyID).(uint32)
	newsletter := r.Context().Value(validator.ResourceKeyID).(*newsletter.Newsletter)
	newsletter.EditorId = userId

	result, err := a.repository.CreateNewsletter(newsletter)
	if err != nil {
		a.logger.Error().Err(err).Msg("unable to create newsletter")
		e.ServerError(w, e.DataCreationFailure)
		return
	}

	a.logger.Info().Str("id", fmt.Sprintf("%d", result.ID)).Msg("newsletter created")
	w.WriteHeader(http.StatusCreated)
}

// Read godoc
//
//	@summary		Read newsletter
//	@description	Read newsletter
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

	userId := r.Context().Value(validator.UserKeyID).(uint32)

	newsletterId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		a.logger.Error().Err(err).Msg("invalid newsletter id in url param")
		e.BadRequest(w, e.InvalidIdInUrlParam)
		return
	}

	newsletter, err2 := a.repository.ReadNewsletter(uint32(newsletterId), userId)
	if err2 != nil {
		if err2 == gorm.ErrRecordNotFound {
			a.logger.Error().Err(err2).Msg("newsletter not found")
			e.NotFoundError(w, e.ResourceNotFound)
			return
		}

		a.logger.Error().Err(err2).Msg("unable to access newsletter")
		e.ServerError(w, e.DataAccessFailure)
		return
	}

	if err := json.NewEncoder(w).Encode(newsletter); err != nil {
		a.logger.Error().Err(err).Msg("unable to encode newsletter into response")
		e.ServerError(w, e.JsonEncodingFailure)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Update godoc
//
//	@summary		Update newsletter
//	@description	Update newsletter
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

	userId := r.Context().Value(validator.UserKeyID).(uint32)
	newsletter := r.Context().Value(validator.ResourceKeyID).(*newsletter.Newsletter)

	newsletterId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		a.logger.Error().Err(err).Msg("invalid newsletter id in url param")
		e.BadRequest(w, e.InvalidIdInUrlParam)
		return
	}

	newsletter.EditorId = userId
	newsletter.ID = uint32(newsletterId)

	if err := a.repository.UpdateNewsletter(newsletter); err != nil {
		if err == gorm.ErrRecordNotFound {
			a.logger.Error().Err(err).Msg("newsletter not found")
			e.NotFoundError(w, e.ResourceNotFound)
			return
		}
		a.logger.Error().Err(err).Msg("unable to update newsletter")
		e.ServerError(w, e.DataUpdateFailure)
		return
	}

	a.logger.Info().Str("id", fmt.Sprintf("%d", newsletter.ID)).Msg("newsletter updated")
	w.WriteHeader(http.StatusOK)
}

// Delete godoc
//
//	@summary		Delete newsletter
//	@description	Delete newsletter
//	@tags			newsletters
//	@accept			json
//	@produce		json
//	@param			id		path		string	true	"Newsletter ID"
//	@success		204
//	@failure		400		{object}	err.Error
//	@failure		500		{object}	err.Error
//	@router			/newsletters/{id} [delete]
func (a *API) Delete(w http.ResponseWriter, r *http.Request) {

	userId := r.Context().Value(validator.UserKeyID).(uint32)

	newsletterId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		a.logger.Error().Err(err).Msg("invalid newsletter id in url param")
		e.BadRequest(w, e.InvalidIdInUrlParam)
		return
	}

	if err := a.repository.DeleteNewsletter(uint32(newsletterId), userId); err != nil {
		if err == gorm.ErrRecordNotFound {
			a.logger.Error().Err(err).Msg("newsletter not found")
			e.NotFoundError(w, e.ResourceNotFound)
			return
		}
		a.logger.Error().Err(err).Msg("unable to delete newsletter")
		e.ServerError(w, e.DataDeletionFailure)
		return
	}

	a.logger.Info().Str("id", fmt.Sprintf("%d", newsletterId)).Msg("newsletter deleted")
	w.WriteHeader(http.StatusNoContent)
}
