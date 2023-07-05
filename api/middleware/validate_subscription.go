package middleware

import (
	"net/http"
	"context"
	"encoding/json"

	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/api/resource/subscription"
	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/util/validator"
	e "github.com/bwc00/strv-go-newsletter-shakleya-mohammed/util/err"
)


// ValidateSubscription is a middleware that validates the subscription entity in the request body.
// It decodes the JSON body into a subscription struct, validates it using the validator, and sets it in the request context.
// If validation fails, it returns the validation errors as the response.
func (m *Middleware) ValidateSubscription(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Decode the JSON body into a subscription struct
		subscription := &subscription.Subscription{}
		if err := json.NewDecoder(r.Body).Decode(subscription); err != nil {
			m.logger.Error().Err(err).Msg("poorly formatted subscription body")
			e.BadRequest(w, e.JsonDecodingFailure)
			return
		}

		// Validate the subscription using the validator
		if err := m.validator.Struct(subscription); err != nil {
			// Convert the validation errors to the error response format
			resp := validator.ToErrResponse(err)
			if resp == nil {
				m.logger.Error().Err(err).Msg("form error response failure")
				e.ServerError(w, e.FormErrResponseFailure)
				return
			}

			// Encode the validation errors as JSON
			respBody, err := json.Marshal(resp)
			if err != nil {
				m.logger.Error().Err(err).Msg("json response encoding failure")
				e.ServerError(w, e.JsonEncodingFailure)
				return
			}

			// Return the validation errors as the response
			e.ValidationErrors(w, respBody)
			return
		}

		// Set the validated subscription in the request context
		ctx := context.WithValue(r.Context(), validator.ResourceKeyID, subscription)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}