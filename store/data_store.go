package store

import (
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
)

type DataStore interface {
	GetPlanIdFromName(planName string) (int, error)
	InsertClientData(planId int, payload types.SignupPayload, accessKey, secretKeyHash string) error
	GetClientFromAccessKey(accessKey string) (*types.ClientData, error)
	InsertUploadMetaData(uploadMetaData *types.UploadMetaData) error
	GetMetaDataByUUID(imgUuid string) (*types.UploadMetaData, error)
	InsertFaceMatchResult(result *types.FaceMatchData) error
	InsertOCRResult(result *types.OCRData) error
	InsertFaceMatchJobCompleted(img1ID, img2ID, clientID int, jobID string) error
	InsertOCRJobCompleted(imgId, clientID int, jobID string) error
	UpdateFaceMatchJobCompleted(jobID string, score int) error
	UpdateOCRJobCompleted(jobID string, result *types.OCRResponse) error
	UpdateFaceMatchJobProcessed(jobID string) error
	UpdateOCRJobProcessed(jobID string) error
	UpdateFaceMatchJobFailed(jobID, reason string) error
	UpdateOCRJobFailed(jobID, reason string) error
}
