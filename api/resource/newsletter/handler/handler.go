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
//	@success		200	{array}		Newsletter
//	@failure		500	{object}	err.Error
//	@router			/newsletters [get]
func (a *API) List(w http.ResponseWriter, r *http.Request) {

	userId := r.Context().Value(validator.KeyID).(uint32)

	newsletters, err := a.repository.ListNewsletters(userId)
	if err != nil {
		a.logger.Error().Err(err).Msg("")
		e.ServerError(w, e.DataAccessFailure)
		return
	}

	if newsletters == nil {
		fmt.Fprint(w, "[]")
		return
	}

	if err := json.NewEncoder(w).Encode(newsletters); err != nil {
		a.logger.Error().Err(err).Msg("")
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
//	@param			name    description
//	@success		201
//	@failure		400	{object}	err.Error
//	@failure		422	{object}	err.Errors
//	@failure		500	{object}	err.Error
//	@router			/newsletters [post]
func (a *API) Create(w http.ResponseWriter, r *http.Request) {

	userId := r.Context().Value(validator.KeyID).(uint32)
	newsletter := r.Context().Value("newsletter").(*newsletter.Newsletter)
	newsletter.EditorId = userId

	result, err := a.repository.CreateNewsletter(newsletter)
	if err != nil {
		a.logger.Error().Err(err).Msg("")
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
//	@param			id
//	@success		200	{object}	Newsletter
//	@failure		400	{object}	err.Error
//	@failure		404
//	@failure		500	{object}	err.Error
//	@router			/newsletters/{id} [get]
func (a *API) Read(w http.ResponseWriter, r *http.Request) {

	userId := r.Context().Value(validator.KeyID).(uint32)

	newsletterId, err2 := strconv.Atoi(chi.URLParam(r, "id"))
	if err2 != nil {
		e.BadRequest(w, e.InvalidIdInUrlParam)
		return
	}

	newsletter, err3 := a.repository.ReadNewsletter(uint32(newsletterId), userId)
	if err3 != nil {
		if err3 == gorm.ErrRecordNotFound {
			a.logger.Error().Err(err3).Msg("")
			e.NotFoundErrors(w, e.ResourceNotFound)
			return
		}

		a.logger.Error().Err(err3).Msg("")
		e.ServerError(w, e.DataAccessFailure)
		return
	}

	if err := json.NewEncoder(w).Encode(newsletter); err != nil {
		a.logger.Error().Err(err).Msg("")
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
//	@param			id		name	description
//	@success		200
//	@failure		400	{object}	err.Error
//	@failure		404
//	@failure		422	{object}	err.Errors
//	@failure		500	{object}	err.Error
//	@router			/newsletters/{id} [put]
func (a *API) Update(w http.ResponseWriter, r *http.Request) {

	userId := r.Context().Value(validator.KeyID).(uint32)
	newsletter := r.Context().Value("newsletter").(*newsletter.Newsletter)

	newsletterId, err2 := strconv.Atoi(chi.URLParam(r, "id"))
	if err2 != nil {
		a.logger.Error().Err(err2).Msg("")
		e.BadRequest(w, e.InvalidIdInUrlParam)
		return
	}

	newsletter.EditorId = userId
	newsletter.ID = uint32(newsletterId)

	if err := a.repository.UpdateNewsletter(newsletter); err != nil {
		if err == gorm.ErrRecordNotFound {
			a.logger.Error().Err(err).Msg("")
			e.NotFoundErrors(w, e.ResourceNotFound)
			return
		}

		a.logger.Error().Err(err).Msg("")
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
//	@param			id
//	@success		200
//	@failure		400	{object}	err.Error
//	@failure		404
//	@failure		500	{object}	err.Error
//	@router			/newsletters/{id} [delete]
func (a *API) Delete(w http.ResponseWriter, r *http.Request) {

	userId := r.Context().Value(validator.KeyID).(uint32)

	newsletterId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		a.logger.Error().Err(err).Msg("")
		e.BadRequest(w, e.InvalidIdInUrlParam)
		return
	}

	if err := a.repository.DeleteNewsletter(uint32(newsletterId), userId); err != nil {
		a.logger.Error().Err(err).Msg("")
		e.ServerError(w, e.DataDeletionFailure)
		return
	}

	a.logger.Info().Str("id", fmt.Sprintf("%d", newsletterId)).Msg("newsletter deleted")
	w.WriteHeader(http.StatusNoContent)
}
