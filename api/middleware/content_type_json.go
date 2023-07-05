package middleware

import "net/http"

// ContentTypeJson is a middleware that sets the Content-Type header to "application/json".
// It wraps the provided `next` http.Handler, allowing it to handle the request and response.
func (m *Middleware) ContentTypeJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the Content-Type header to JSON
		w.Header().Set("Content-Type", "application/json;charset=utf8")

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
