package service

import (
	"mime/multipart"

	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
)

type ServiceManager interface {
	SignupClient(payload types.SignupPayload) (*KeyPair, error)
	ValidateFile(fileName, fileType string) error
	SaveFile(fileHeader *multipart.FileHeader, uploadMetaData *types.UploadMetaData) error
	PerformFaceMatch(payload types.FaceMatchPayload, clientID int) (string, error)
	PerformOCR(payload types.OCRPayload, clientID int) (string, error)
	GetJobDetailsByJobID(jobID, jobType string) (*types.JobRecord, error)
	FetchDataFromCache(payload interface{}, clientID int, jobType string) (string, bool)
	SetDataInCache(payload interface{}, clientID int, jobType, jobID string)
}
