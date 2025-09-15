package config

import (
	"github.com/go-playground/validator/v10"
	// "github.com/iyasz/golang-clean-architecture/internal/config"
)

func NewValidator() *validator.Validate {
	return validator.New()
}
