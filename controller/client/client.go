package client

import (
	"errors"
	"regexp"
	"time"

	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
)

var ErrInvalidEmail = errors.New("invalid email")
var ErrInvalidPlan = errors.New("invalid plan, supported plans are basic, advanced, or enterprise")

const AccessTokenExpiry = 15 * time.Minute
const RefreshTokenExpiry = 7 * 24 * time.Hour

// TODO: Change the interface name
type ClientServiceInterface interface {
	ValidatePayload(payload types.SignupPayload) error
	GenerateTokenPair(payload types.SignupPayload) (TokenPair, error)
}

type ClientService struct {
	tokenService TokenGenerator
}

func NewClientService(tokenService TokenGenerator) ClientService {
	return ClientService{
		tokenService: tokenService,
	}
}

func (c ClientService) ValidatePayload(payload types.SignupPayload) error {
	if err := validateEmail(payload.Email); err != nil {
		return err
	}
	if err := validatePlan(payload.Plan); err != nil {
		return err
	}

	return nil
}

func (c ClientService) GenerateTokenPair(payload types.SignupPayload) (*TokenPair, error) {
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
