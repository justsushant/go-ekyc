package service

import (
	"mime/multipart"

	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
)

type ServiceManager interface {
	SignupClient(payload types.SignupPayload) (*KeyPair, error)
	ValidateFile(fileName, fileType string) error
	SaveFile(fileHeader *multipart.FileHeader, uploadMetaData *types.UploadMetaData) error
	ValidateImage(payload types.FaceMatchPayload, clientID int) error
	CalcAndSaveFaceMatchScore(payload types.FaceMatchPayload, clientID int) (int, error)
	ValidateImageOCR(payload types.OCRPayload, clientID int) error
	PerformAndSaveOCR(payload types.OCRPayload, clientID int) (*types.OCRResponse, error)
	PerformFaceMatchAsync(payload types.FaceMatchPayload, clientID int) (string, error)
	PerformOCRAsync(payload types.OCRPayload, clientID int) (string, error)
	GetJobDetailsByJobID(jobID, jobType string) (*types.JobRecord, error)
	FetchDataFromCache(payload interface{}, clientID int, jobType string) (string, bool)
	SetDataInCache(payload interface{}, clientID int, jobType, jobID string)
}
