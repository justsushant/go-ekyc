package store

import "github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"

type WorkerDataStore interface {
	UpdateFaceMatchJobCompleted(jobID string, score int) error
	UpdateOCRJobCompleted(jobID string, result *types.OCRResponse) error
	UpdateFaceMatchJobProcessed(jobID string) error
	UpdateOCRJobProcessed(jobID string) error
	UpdateFaceMatchJobFailed(jobID, reason string) error
	UpdateOCRJobFailed(jobID, reason string) error
}
