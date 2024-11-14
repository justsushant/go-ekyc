package controller

import "github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"

// TODO: Change the interface name
type ControllerInterface interface {
	ValidatePayload(payload types.SignupPayload) error
	GenerateTokenPair(payload types.SignupPayload) (*TokenPair, error)
}
