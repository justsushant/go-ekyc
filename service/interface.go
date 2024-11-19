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
	SaveFile(fileHeader *multipart.FileHeader, uploadMetaData *types.UploadMetaData) error
	ValidateImage(payload types.FaceMatchPayload, clientID int) error
	CalcAndSaveFaceMatchScore(payload types.FaceMatchPayload, clientID int) (int, error)
	ValidateImageOCR(payload types.OCRPayload, clientID int) error
	PerformAndSaveOCR(payload types.OCRPayload, clientID int) (*types.OCRResponse, error)
}
