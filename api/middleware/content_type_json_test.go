package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/api/middleware"
	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/util/logger"
	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/util/validator"
)

var (
	expRespBody    = "{\"message\":\"Hello World!\"}"
	expContentType = "application/json;charset=utf8"
)

// TestContentTypeJson tests the ContentTypeJson middleware.
func TestContentTypeJson(t *testing.T) {

	// Create a new HTTP request and response recorder
	r, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	// Initialize the middleware handlers
	middlewareHandlers := middleware.New(logger.New(true), validator.New())

	// Invoke the ContentTypeJson middleware
	middlewareHandlers.ContentTypeJson(http.HandlerFunc(sampleHandlerFunc())).ServeHTTP(rr, r)
	response := rr.Result()

	// Check the response body
	if respBody := rr.Body.String(); respBody != expRespBody {
		t.Errorf("Wrong response body:  got %v want %v ", respBody, expRespBody)
	}

	// Check the status code
	if status := response.StatusCode; status != http.StatusOK {
		t.Errorf("Wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the Content-Type header
	if contentType := response.Header.Get("Content-type"); contentType != expContentType {
		t.Errorf("Wrong status code: got %v want %v", contentType, expContentType)
	}
}

// sampleHandlerFunc returns a sample handler function that writes the expected response body.
func sampleHandlerFunc() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expRespBody)
	}
}
