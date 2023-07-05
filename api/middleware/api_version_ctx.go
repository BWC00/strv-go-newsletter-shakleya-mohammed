package middleware

import (
	"net/http"
	"context"

	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/util/validator"
)

func (m *Middleware) ApiVersionCtx(version string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), validator.ApiVersionKeyID, version))
			next.ServeHTTP(w, r)
		})
	}
}