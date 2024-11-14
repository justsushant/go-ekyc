package controller

import (
	"errors"
)

var (
	ErrInvalidEmail = errors.New("invalid email")
	ErrInvalidPlan  = errors.New("invalid plan, supported plans are basic, advanced, or enterprise")
)
