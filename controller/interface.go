package controller

import (
	"mime/multipart"

	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
)

// TODO: Change the interface name
type ControllerInterface interface {
	ValidatePayload(payload types.SignupPayload) error
	GenerateTokenPair(payload types.SignupPayload) (*TokenPair, error)
	SaveSignupData(payload types.SignupPayload, refreshToken string) error
	ValidateFile(fileName, fileType string) error
	SaveUploadedFile(fileHeader *multipart.FileHeader) error
}

type Store interface {
	GetPlanIdFromName(planName string) (int, error)
	InsertClientData(payload types.SignupPayload, planId int, refreshToken string) error
}

type FileStore interface {
	SaveFile(fileHeader *multipart.FileHeader) error
}
