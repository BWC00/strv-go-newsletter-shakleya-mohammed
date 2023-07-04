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
			m.logger.Error().Err(err).Msg("")
			e.ServerError(w, e.UnauthorizedAccess)
			return
		}

		userId, err := auth.ExtractTokenID(r)
		if err != nil {
			m.logger.Error().Err(err).Msg("")
			e.ServerError(w, e.UnauthorizedAccess)
		}

		ctx := context.WithValue(r.Context(), validator.KeyID, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}