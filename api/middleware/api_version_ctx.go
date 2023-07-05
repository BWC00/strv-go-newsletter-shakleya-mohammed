package middleware

import (
	"net/http"
	"context"

	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/util/validator"
)

// ApiVersionCtx sets the API version in the request context.
// It returns a middleware handler that wraps the provided `next` http.Handler.
func (m *Middleware) ApiVersionCtx(version string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set the API version in the request context
			r = r.WithContext(context.WithValue(r.Context(), validator.ApiVersionKeyID, version))

			// Call the next handler in the chain
			next.ServeHTTP(w, r)
		})
	}
}