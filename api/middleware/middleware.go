package middleware

import (
	vd "github.com/go-playground/validator/v10"

	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/util/logger"
)

// Middleware represents the collection of middleware handlers.
type Middleware struct {
	logger     *logger.Logger
	validator  *vd.Validate
}

// New creates a new instance of the Middleware struct.
// It takes a logger and a validator as parameters and returns a pointer to the Middleware struct.
func New(logger *logger.Logger, validator *vd.Validate) *Middleware {
	return &Middleware{
		logger:     logger,
		validator:  validator,
	}
}