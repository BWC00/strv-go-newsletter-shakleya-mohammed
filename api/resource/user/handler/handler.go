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

// API represents the handler for the user API endpoints.
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

// Register godoc
//
//	@summary		Register user
//	@description	Register user
//	@tags			users
//	@accept			json
//	@produce		json
//	@param			body	body	user.User	true	"User contents"
//	@success		201
//	@failure		400	{object}	err.Error
//	@failure		422	{object}	err.Errors
//	@failure		500	{object}	err.Error
//	@router			/register [post]
func (a *API) Register(w http.ResponseWriter, r *http.Request) {
	// Extract user from the request context
	user := r.Context().Value(validator.ResourceKeyID).(*user.User)

	// Register the user
	token, err := a.repository.RegisterUser(user)
	if err != nil {
		// Check if the email is already taken
		if err.Error() == e.FieldNotUnique {
			a.logger.Error().Err(err).Msg("email already taken")
			e.BadRequest(w, e.FieldNotUnique)
			return
		}

		// Handle other errors during user registration
		a.logger.Error().Err(err).Msg("unable to register user")
		e.ServerError(w, e.DataCreationFailure)
		return
	}

	// Encode the token into the response
	if err := json.NewEncoder(w).Encode(token); err != nil {
		a.logger.Error().Err(err).Msg("unable to encode token in response")
		e.ServerError(w, e.JsonEncodingFailure)
		return
	}

	// Log successful user creation
	a.logger.Info().Str("id", fmt.Sprintf("%d", user.ID)).Msg("new user created")

	// Return 201 Created status
	w.WriteHeader(http.StatusCreated)
}


// Login godoc
//
//	@summary		Login user
//	@description	Login user
//	@tags			users
//	@accept			json
//	@produce		json
//	@param			body	body	user.User	true	"User contents"
//	@success		200
//	@failure		400	{object}	err.Error
//	@failure		422	{object}	err.Errors
//	@failure		500	{object}	err.Error
//	@router			/login [post]
func (a *API) Login(w http.ResponseWriter, r *http.Request) {
	// Extract user from the request context
	user := r.Context().Value(validator.ResourceKeyID).(*user.User)

	// Login the user
	token, err := a.repository.LoginUser(user)
	if err != nil {
		// Check if the email doesn't exist
		if err == gorm.ErrRecordNotFound {
			a.logger.Error().Err(err).Msg("email doesn't exist")
			e.BadRequest(w, e.ResourceNotFound)
			return
		}

		// Handle incorrect password error
		a.logger.Error().Err(err).Msg("password not correct")
		e.BadRequest(w, e.AuthenticationFailure)
		return
	}

	// Encode the token into the response
	if err := json.NewEncoder(w).Encode(token); err != nil {
		a.logger.Error().Err(err).Msg("unable to encode token in response")
		e.ServerError(w, e.JsonEncodingFailure)
		return
	}

	// Log successful user login
	a.logger.Info().Str("id", fmt.Sprintf("%d", user.ID)).Msg("user logged in")

	// Return 200 OK status
	w.WriteHeader(http.StatusOK)
}