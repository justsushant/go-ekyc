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
	InsertFaceMatchJobCreated(img1ID, img2ID, clientID int, jobID string) error
	InsertOCRJobCreated(imgId, clientID int, jobID string) error
	GetFaceMatchByJobID(jobID string) (*types.JobRecord, error)
	GetOCRByJobID(jobID string) (*types.JobRecord, error)
}
