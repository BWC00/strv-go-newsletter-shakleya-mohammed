package middleware

import (
	vd "github.com/go-playground/validator/v10"

	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/util/logger"
)

type Middleware struct {
	logger     *logger.Logger
	validator  *vd.Validate
}

func New(logger *logger.Logger, validator *vd.Validate) *Middleware {
	return &Middleware{
		logger:     logger,
		validator:  validator,
	}
}