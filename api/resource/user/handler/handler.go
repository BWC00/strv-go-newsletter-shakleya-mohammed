package handler

import (
	"encoding/json"
	"net/http"
	"fmt"

	vd "github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/api/resource/user/repository"
	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/api/resource/user"
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

// Register godoc
//
//	@summary		Register user
//	@description	Register user
//	@tags			users
//	@accept			json
//	@produce		json
//	@param			body	body	User	true	"User contents"
//	@success		201
//	@failure		400	{object}	err.Error
//	@failure		422	{object}	err.Errors
//	@failure		500	{object}	err.Error
//	@router			/register [post]
func (a *API) Register(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value(validator.ResourceKeyID).(*user.User)

	token, err := a.repository.RegisterUser(user)
	if err != nil {
		if err.Error() == e.FieldNotUnique {
			a.logger.Error().Err(err).Msg("email already taken")
			e.BadRequest(w, e.FieldNotUnique)
			return
		}
		a.logger.Error().Err(err).Msg("unable to register user")
		e.ServerError(w, e.DataCreationFailure)
		return
	}

	if err := json.NewEncoder(w).Encode(token); err != nil {
		a.logger.Error().Err(err).Msg("unable to encode token in response")
		e.ServerError(w, e.JsonEncodingFailure)
		return
	}

	a.logger.Info().Str("id", fmt.Sprintf("%d", user.ID)).Msg("new user created")
	w.WriteHeader(http.StatusCreated)
}


// Login godoc
//
//	@summary		Login user
//	@description	Login user
//	@tags			users
//	@accept			json
//	@produce		json
//	@param			body	body	User	true	"User contents"
//	@success		200
//	@failure		400	{object}	err.Error
//	@failure		422	{object}	err.Errors
//	@failure		500	{object}	err.Error
//	@router			/login [post]
func (a *API) Login(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value(validator.ResourceKeyID).(*user.User)

	token, err := a.repository.LoginUser(user)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			a.logger.Error().Err(err).Msg("email doesn't exist")
			e.BadRequest(w, e.ResourceNotFound)
			return
		}
		a.logger.Error().Err(err).Msg("password not correct")
		e.BadRequest(w, e.AuthenticationFailure)
		return
	}

	if err := json.NewEncoder(w).Encode(token); err != nil {
		a.logger.Error().Err(err).Msg("unable to encode token in response")
		e.ServerError(w, e.JsonEncodingFailure)
		return
	}

	a.logger.Info().Str("id", fmt.Sprintf("%d", user.ID)).Msg("user logged in")
}