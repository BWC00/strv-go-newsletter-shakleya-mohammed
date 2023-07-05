package middleware

import (
	"net/http"
	"context"
	
	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/util/validator"
	e "github.com/bwc00/strv-go-newsletter-shakleya-mohammed/util/err"
	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/util/auth"
)

// Authentication is a middleware that performs authentication for incoming requests.
// It checks the validity of the token in the request header, extracts the user ID
// from the token, and sets it in the request context before calling the next handler in the chain.
func (m *Middleware) Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the validity of the token
		if err := auth.TokenValid(r); err != nil {
			m.logger.Error().Err(err).Msg("token invalid")
			e.ServerError(w, e.UnauthorizedAccess)
			return
		}

		// Extract the user ID from the token
		userId, err := auth.ExtractTokenID(r)
		if err != nil {
			m.logger.Error().Err(err).Msg("couldn't extract userID from token")
			e.ServerError(w, e.TokenExtractionFailure)
		}

		// Set the user ID in the request context
		ctx := context.WithValue(r.Context(), validator.UserKeyID, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}