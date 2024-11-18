package service

import (
	"mime/multipart"

	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
)

// TODO: Change the interface name
type ControllerInterface interface {
	ValidatePayload(payload types.SignupPayload) error
	GenerateKeyPair() (*KeyPair, error)
	SaveSignupData(payload types.SignupPayload, keyPair *KeyPair) error
	ValidateFile(fileName, fileType string) error
	SaveUploadedFile(fileHeader *multipart.FileHeader) error
}
