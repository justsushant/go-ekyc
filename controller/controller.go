package controller

import (
	"regexp"

	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
)

type Service struct {
	tokenService TokenGenerator
}

func NewService(tokenService TokenGenerator) Service {
	return Service{
		tokenService: tokenService,
	}
}

func (c Service) ValidatePayload(payload types.SignupPayload) error {
	if err := validateEmail(payload.Email); err != nil {
		return err
	}
	if err := validatePlan(payload.Plan); err != nil {
		return err
	}

	return nil
}

func (c Service) GenerateTokenPair(payload types.SignupPayload) (*TokenPair, error) {
	return c.tokenService.GenerateTokenPair(payload)
}

func validatePlan(plan string) error {
	switch plan {
	case types.BasicPlan, types.AdvancePlan, types.EnterprisePlan:
		return nil
	default:
		return ErrInvalidPlan
	}
}

func validateEmail(email string) error {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	if re.MatchString(email) {
		return nil
	}
	return ErrInvalidEmail
}
