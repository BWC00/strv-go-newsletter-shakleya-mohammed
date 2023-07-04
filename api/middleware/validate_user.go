package middleware

import (
	"net/http"
	"context"
	"encoding/json"

	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/api/resource/user"
	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/util/validator"
	e "github.com/bwc00/strv-go-newsletter-shakleya-mohammed/util/err"
)

func (m *Middleware) ValidateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user := &user.User{}
		if err := json.NewDecoder(r.Body).Decode(user); err != nil {
			m.logger.Error().Err(err).Msg("")
			e.BadRequest(w, e.JsonDecodingFailure)
			return
		}

		if err := m.validator.Struct(user); err != nil {
			resp := validator.ToErrResponse(err)
			if resp == nil {
				e.ServerError(w, e.FormErrResponseFailure)
				return
			}

			respBody, err := json.Marshal(resp)
			if err != nil {
				m.logger.Error().Err(err).Msg("")
				e.ServerError(w, e.JsonEncodingFailure)
				return
			}

			e.ValidationErrors(w, respBody)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}