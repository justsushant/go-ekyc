package worker

import (
	"database/sql"

	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/db"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
)

type PsqlWorkerStore struct {
	db *sql.DB
}

func NewPsqlWorkerStore(dsn string) PsqlWorkerStore {
	psqlClient := db.NewPsqlClient(dsn)
	return PsqlWorkerStore{
		db: psqlClient,
	}
}

func (s PsqlWorkerStore) UpdateFaceMatchJobCompleted(jobID string, score int) error {
	_, err := s.db.Exec(
		"UPDATE face_match SET match_score = $1, completed_at = NOW(), status = $2 WHERE job_id = $3",
		score, types.JobStatusCompleted, jobID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s PsqlWorkerStore) UpdateOCRJobCompleted(jobID string, data *types.OCRResponse) error {
	_, err := s.db.Exec(
		"UPDATE ocr SET details = $1, completed_at = NOW(), status = $2 WHERE job_id = $3",
		data.String(), types.JobStatusCompleted, jobID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s PsqlWorkerStore) UpdateFaceMatchJobProcessed(jobID string) error {
	_, err := s.db.Exec(
		"UPDATE face_match SET processed_at = NOW(), status = $1 WHERE job_id = $2",
		types.JobStatusProcessing, jobID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s PsqlWorkerStore) UpdateOCRJobProcessed(jobID string) error {
	_, err := s.db.Exec(
		"UPDATE ocr SET processed_at = NOW(), status = $1 WHERE job_id = $2",
		types.JobStatusProcessing, jobID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s PsqlWorkerStore) UpdateFaceMatchJobFailed(jobID, reason string) error {
	_, err := s.db.Exec(
		"UPDATE face_match SET failed_at = NOW(), status = $1, failed_reason = $2 WHERE job_id = $3",
		types.JobStatusFailed, reason, jobID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s PsqlWorkerStore) UpdateOCRJobFailed(jobID, reason string) error {
	_, err := s.db.Exec(
		"UPDATE ocr SET failed_at = NOW(), status = $1, failed_reason = $2 WHERE job_id = $3",
		types.JobStatusFailed, reason, jobID,
	)
	if err != nil {
		return err
	}

	return nil
}
