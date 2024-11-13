package client

import (
	"errors"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
)

var ErrInvalidEmail = errors.New("invalid email")
var ErrInvalidPlan = errors.New("invalid plan, supported plans are basic, advanced, or enterprise")

const AccessTokenExpiry = 15 * time.Minute
const RefreshTokenExpiry = 7 * 24 * time.Hour

// TODO: Change the interface name
type ClientServiceInterface interface {
	ValidatePayload(payload types.SignupPayload) error
	GenerateAccessToken(payload types.SignupPayload, expiryTime time.Duration, secret []byte) (string, error)
	GenerateRefreshToken(payload types.SignupPayload, expiryTime time.Duration, secret []byte) (string, error)
}

type ClientService struct{}

func NewClientService() ClientService {
	return ClientService{}
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

func (c ClientService) GenerateAccessToken(payload types.SignupPayload, expiryTime time.Duration, secret []byte) (string, error) {
	claims := jwt.MapClaims{
		"client_email": payload.Email,
		"client_plan":  payload.Plan,
		"exp":          time.Now().Add(expiryTime).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func (c ClientService) GenerateRefreshToken(payload types.SignupPayload, expiryTime time.Duration, secret []byte) (string, error) {
	claims := jwt.MapClaims{
		"client_email": payload.Email,
		"client_plan":  payload.Plan,
		"exp":          time.Now().Add(expiryTime).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
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
