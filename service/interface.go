package service

import (
	"mime/multipart"

	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
)

// TODO: Change the interface name
type ControllerInterface interface {
	ValidatePayload(payload types.SignupPayload) error
	GenerateKeyPair(payload types.SignupPayload) (*KeyPair, error)
	SaveSignupData(payload types.SignupPayload, refreshToken string) error
	ValidateFile(fileName, fileType string) error
	SaveUploadedFile(fileHeader *multipart.FileHeader) error
}
