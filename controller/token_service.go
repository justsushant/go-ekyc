package controller

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
)

const ACCESS_TOKEN_EXPIRY = 15 * time.Minute
const REFRESH_TOKEN_EXPIRY = 7 * 24 * time.Hour

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

// TODO: No idea why its throwing error (no .env file found) when screts are accessed using config from inside
func NewTokenService(accessKeySecret, refreshKeySecret string) TokenService {
	return TokenService{
		AccessTokenExpiry:  ACCESS_TOKEN_EXPIRY,
		RefreshTokenExpiry: REFRESH_TOKEN_EXPIRY,
		AccessTokenSecret:  []byte(accessKeySecret),
		RefreshTokenSecret: []byte(refreshKeySecret),
	}
}

// func NewTokenService() TokenService {
// 	return TokenService{
// 		AccessTokenExpiry:  ACCESS_TOKEN_EXPIRY,
// 		RefreshTokenExpiry: REFRESH_TOKEN_EXPIRY,
// 		AccessTokenSecret:  []byte(config.Envs.Access_token_secret),
// 		RefreshTokenSecret: []byte(config.Envs.Refresh_token_secret),
// 	}
// }

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
