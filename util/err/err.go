package err

import (
	"fmt"
	"net/http"
)

// Constants for common error messages
const (
	DataCreationFailure = "data creation failure"
	DataAccessFailure   = "data access failure"
	DataUpdateFailure   = "data update failure"
	DataDeletionFailure = "data deletion failure"

	JsonEncodingFailure = "json encoding failure"
	JsonDecodingFailure = "json decoding failure"

	FormErrResponseFailure = "form error response failure"
	InvalidIdInUrlParam = "invalid id in url param"
	UnauthorizedAccess = "unauthorized access"
	ResourceNotFound = "resource not found"
	SendingEmailFailure = "sending email failure"
	AuthenticationFailure = "authentication failure"
	TokenExtractionFailure = "token extraction failure"
	FieldNotUnique = "email not unique"

	EndpointNotFound = "endpoint not found"

)

// Error represents a single error message.
type Error struct {
	Error string `json:"error"`
}

// Errors represents a collection of error messages.
type Errors struct {
	Errors []string `json:"errors"`
}

// ServerError sends a HTTP 500 Internal Server Error response with the specified error message.
func ServerError(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, `{"error": "%v"}`, msg)
}

// BadRequest sends a HTTP 400 Bad Request response with the specified error message.
func BadRequest(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, `{"error": "%v"}`, msg)
}

// ValidationErrors sends a HTTP 422 Unprocessable Entity response with the specified error representation.
func ValidationErrors(w http.ResponseWriter, reps []byte) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	w.Write(reps)
}

// Unauthorized sends a HTTP 401 Unauthorized response with the specified error message.
func Unauthorized(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusUnauthorized)
	fmt.Fprintf(w, `{"error": "%v"}`, msg)
}

// NotFoundError sends a HTTP 404 Not Found response with the specified error message.
func NotFoundError(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, `{"error": "%v"}`, msg)
}
