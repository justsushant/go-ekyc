package client

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type TokenGenerator interface {
	GenerateTokenPair(payload types.SignupPayload) (*TokenPair, error)
}

type TokenService struct {
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	AccessTokenSecret  []byte
	RefreshTokenSecret []byte
}

func NewTokenService(accessTokenExpiry, refreshTokenExpiry time.Duration, accessTokenSecret, refreshTokenSecret []byte) TokenService {
	return TokenService{
		AccessTokenExpiry:  accessTokenExpiry,
		RefreshTokenExpiry: refreshTokenExpiry,
		AccessTokenSecret:  accessTokenSecret,
		RefreshTokenSecret: refreshTokenSecret,
	}
}

// TODO: Consider security and token length considerations
func (t TokenService) GenerateTokenPair(payload types.SignupPayload) (*TokenPair, error) {
	accessToken, err := t.generateToken(payload, t.AccessTokenExpiry, t.AccessTokenSecret)
	if err != nil {
		return nil, err
	}

	refreshToken, err := t.generateToken(payload, t.RefreshTokenExpiry, t.RefreshTokenSecret)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (t TokenService) generateToken(payload types.SignupPayload, expiryTime time.Duration, secret []byte) (string, error) {
	claims := jwt.MapClaims{
		"client_email": payload.Email,
		"client_plan":  payload.Plan,
		"exp":          time.Now().Add(expiryTime).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
