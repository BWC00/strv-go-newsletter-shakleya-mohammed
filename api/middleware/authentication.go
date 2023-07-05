package middleware

import (
	"net/http"
	"context"
	
	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/util/validator"
	e "github.com/bwc00/strv-go-newsletter-shakleya-mohammed/util/err"
	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/util/auth"
)

func (m *Middleware) Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := auth.TokenValid(r); err != nil {
			m.logger.Error().Err(err).Msg("token invalid")
			e.ServerError(w, e.UnauthorizedAccess)
			return
		}

		userId, err := auth.ExtractTokenID(r)
		if err != nil {
			m.logger.Error().Err(err).Msg("couldn't extract userID from token")
			e.ServerError(w, e.TokenExtractionFailure)
		}

		ctx := context.WithValue(r.Context(), validator.UserKeyID, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}