package middleware_test

// import (
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/api/middleware"
// 	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/config"
// 	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/util/logger"
// 	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/util/validator"
// )

// var (
// 	expRespBody    = "{\"message\":\"Hello World!\"}"
// 	expContentType = "application/json;charset=utf8"
// )

// func TestContentTypeJson(t *testing.T) {
// 	r, _ := http.NewRequest("GET", "/", nil)
// 	rr := httptest.NewRecorder()

// 	cfg := config.New()

// 	middlewareHandlers := middleware.New(logger.New(cfg.Server.Debug), validator.New())

// 	middlewareHandlers.ContentTypeJson(http.HandlerFunc(sampleHandlerFunc())).ServeHTTP(rr, r)
// 	response := rr.Result()

// 	if respBody := rr.Body.String(); respBody != expRespBody {
// 		t.Errorf("Wrong response body:  got %v want %v ", respBody, expRespBody)
// 	}

// 	if status := response.StatusCode; status != http.StatusOK {
// 		t.Errorf("Wrong status code: got %v want %v", status, http.StatusOK)
// 	}

// 	if contentType := response.Header.Get("Content-type"); contentType != expContentType {
// 		t.Errorf("Wrong status code: got %v want %v", contentType, expContentType)
// 	}
// }

// func sampleHandlerFunc() func(w http.ResponseWriter, r *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Fprint(w, expRespBody)
// 	}
// }
