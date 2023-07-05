package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Custom key type.
type key int

const (
	alphaZeroRegexString string = "^[a-zA-Z]*$"

	UserKeyID key = iota
    ResourceKeyID key = iota
    ApiVersionKeyID key = iota
)

// ErrResponse represents an error response containing a list of error messages.
type ErrResponse struct {
	Errors []string `json:"errors"`
}

// New creates a new validator instance using the go-playground validator library.
// It registers custom validation tags and tag name functions.
func New() *validator.Validate {
	validate := validator.New()
	validate.SetTagName("form")

	// Using the names which have been specified for JSON representations of structs, rather than normal Go field names
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	validate.RegisterValidation("alpha_zero", isAlphaZero)

	return validate
}

// ToErrResponse converts a validation error to an ErrResponse containing formatted error messages.
func ToErrResponse(err error) *ErrResponse {
	if fieldErrors, ok := err.(validator.ValidationErrors); ok {
		resp := ErrResponse{
			Errors: make([]string, len(fieldErrors)),
		}

		for i, err := range fieldErrors {
			switch err.Tag() {
			case "required":
				resp.Errors[i] = fmt.Sprintf("%s is a required field", err.Field())
			case "max":
				resp.Errors[i] = fmt.Sprintf("%s must be a maximum of %s in length", err.Field(), err.Param())
			case "email":
				resp.Errors[i] = fmt.Sprintf("%s must be a valid Email", err.Field())
			case "alpha_zero":
				resp.Errors[i] = fmt.Sprintf("%s can only contain alphabetic and space characters", err.Field())
			default:
				resp.Errors[i] = fmt.Sprintf("something wrong on %s; %s", err.Field(), err.Tag())
			}
		}

		return &resp
	}

	return nil
}

// isAlphaZero is a custom validation function that checks if a string contains only alphabetic characters.
func isAlphaZero(fl validator.FieldLevel) bool {
	reg := regexp.MustCompile(alphaZeroRegexString)
	return reg.MatchString(fl.Field().String())
}
