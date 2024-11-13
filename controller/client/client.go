package client

import (
	"errors"
	"regexp"

	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
)

var ErrInvalidEmail = errors.New("invalid email")
var ErrInvalidPlan = errors.New("invalid plan, supported plans are basic, advanced, or enterprise")

func ValidatePayload(payload types.SignupPayload) error {
	if err := isEmailValid(payload.Email); err != nil {
		return err
	}
	if err := isPlanValid(payload.Plan); err != nil {
		return err
	}

	return nil
}

func isPlanValid(plan string) error {
	switch plan {
	case types.BasicPlan, types.AdvancePlan, types.EnterprisePlan:
		return nil
	default:
		return ErrInvalidPlan
	}
}

func isEmailValid(email string) error {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	if re.MatchString(email) {
		return nil
	}
	return ErrInvalidEmail
}
