package middleware

import (
	"net/http"
	"context"
	"encoding/json"

	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/api/resource/newsletter"
	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/util/validator"
	e "github.com/bwc00/strv-go-newsletter-shakleya-mohammed/util/err"
)

func (m *Middleware) ValidateNewsletter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		newsletter := &newsletter.Newsletter{}
		if err := json.NewDecoder(r.Body).Decode(newsletter); err != nil {
			m.logger.Error().Err(err).Msg("poorly formatted newsletter body")
			e.BadRequest(w, e.JsonDecodingFailure)
			return
		}

		if err := m.validator.Struct(newsletter); err != nil {
			resp := validator.ToErrResponse(err)
			if resp == nil {
				m.logger.Error().Err(err).Msg("form error response failure")
				e.ServerError(w, e.FormErrResponseFailure)
				return
			}

			respBody, err := json.Marshal(resp)
			if err != nil {
				m.logger.Error().Err(err).Msg("json response encoding failure")
				e.ServerError(w, e.JsonEncodingFailure)
				return
			}

			e.ValidationErrors(w, respBody)
			return
		}

		ctx := context.WithValue(r.Context(), validator.KeyID, newsletter)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}