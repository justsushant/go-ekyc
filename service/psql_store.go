package service

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/db"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
)

type PsqlStore struct {
	db *sql.DB
}

func NewPsqlStore(dsn string) PsqlStore {
	psqlClient := db.NewPsqlClient(dsn)
	return PsqlStore{
		db: psqlClient,
	}
}

func (s PsqlStore) InsertClientData(planId int, payload types.SignupPayload, accessKey, secretKeyHash string) error {
	_, err := s.db.Exec(
		"INSERT INTO client (name, email, access_key, secret_key_hash, plan_id) VALUES ($1, $2, $3, $4, $5)",
		payload.Name, payload.Email, accessKey, secretKeyHash, planId,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s PsqlStore) GetPlanIdFromName(planName string) (int, error) {
	var planId int
	err := s.db.QueryRow("SELECT id FROM plan WHERE name = $1", planName).Scan(&planId)
	if err != nil {
		return 0, err
	}

	return planId, nil
}

func (s PsqlStore) GetClientFromAccessKey(accessKey string) (*types.ClientData, error) {
	var clientData types.ClientData
	err := s.db.QueryRow(
		"SELECT id, name, email, plan_id, access_key, secret_key_hash FROM client WHERE access_key = $1",
		accessKey,
	).Scan(
		&clientData.Id,
		&clientData.Name,
		&clientData.Email,
		&clientData.PlanID,
		&clientData.AccessKey,
		&clientData.SecretKeyHash,
	)
	if err != nil {
		return nil, err
	}

	return &clientData, nil
}

func (s PsqlStore) InsertUploadMetaData(uploadMetaData *types.UploadMetaData) error {
	_, err := s.db.Exec(
		"INSERT INTO upload (type, client_id, file_path, file_size_kb) VALUES ($1, $2, $3, $4)",
		uploadMetaData.Type, uploadMetaData.ClientID, uploadMetaData.FilePath, uploadMetaData.FileSizeKB,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s PsqlStore) GetMetaDataByUUID(imgUuid string) (*types.UploadMetaData, error) {
	var uploadData types.UploadMetaData
	err := s.db.QueryRow(
		"SELECT id, type, client_id, file_path, file_size_kb FROM upload WHERE file_path LIKE '%' || $1 || '%'",
		imgUuid,
	).Scan(
		&uploadData.Id,
		&uploadData.Type,
		&uploadData.ClientID,
		&uploadData.FilePath,
		&uploadData.FileSizeKB,
	)
	if err != nil {
		return nil, err
	}

	return &uploadData, nil
}

func (s PsqlStore) InsertFaceMatchResult(result *types.FaceMatchData) error {
	_, err := s.db.Exec(
		"INSERT INTO face_match (client_id, upload_id1, upload_id2, match_score) VALUES ($1, $2, $3, $4)",
		result.ClientID, result.ImageID1, result.ImageID2, result.Score,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s PsqlStore) InsertOCRResult(result *types.OCRData) error {
	_, err := s.db.Exec(
		"INSERT INTO ocr (client_id, upload_id, details) VALUES ($1, $2, $3)",
		result.ClientID, result.ImageID, result.Data,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s PsqlStore) InsertFaceMatchJobCreated(img1ID, img2ID, clientID int, jobID string) error {
	_, err := s.db.Exec(
		"INSERT INTO face_match (job_id, status, client_id, upload_id1, upload_id2) VALUES ($1, $2, $3, $4, $5)",
		jobID, types.JobStatusCreated, clientID, img1ID, img2ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s PsqlStore) InsertOCRJobCreated(imgID, clientID int, jobID string) error {
	_, err := s.db.Exec(
		"INSERT INTO ocr (job_id, status, client_id, upload_id) VALUES ($1, $2, $3, $4)",
		jobID, types.JobStatusCreated, clientID, imgID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s PsqlStore) UpdateFaceMatchJobCompleted(jobID string, score int) error {
	_, err := s.db.Exec(
		"UPDATE face_match SET match_score = $1, completed_at = NOW(), status = $2 WHERE job_id = $3",
		score, types.JobStatusCompleted, jobID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s PsqlStore) UpdateOCRJobCompleted(jobID string, data *types.OCRResponse) error {
	_, err := s.db.Exec(
		"UPDATE ocr SET details = $1, completed_at = NOW(), status = $2 WHERE job_id = $3",
		data.String(), types.JobStatusCompleted, jobID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s PsqlStore) UpdateFaceMatchJobProcessed(jobID string) error {
	_, err := s.db.Exec(
		"UPDATE face_match SET processed_at = NOW(), status = $1 WHERE job_id = $2",
		types.JobStatusProcessing, jobID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s PsqlStore) UpdateOCRJobProcessed(jobID string) error {
	_, err := s.db.Exec(
		"UPDATE ocr SET processed_at = NOW(), status = $1 WHERE job_id = $2",
		types.JobStatusProcessing, jobID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s PsqlStore) UpdateFaceMatchJobFailed(jobID, reason string) error {
	_, err := s.db.Exec(
		"UPDATE face_match SET failed_at = NOW(), status = $1, failed_reason = $2 WHERE job_id = $3",
		types.JobStatusFailed, reason, jobID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s PsqlStore) UpdateOCRJobFailed(jobID, reason string) error {
	_, err := s.db.Exec(
		"UPDATE ocr SET failed_at = NOW(), status = $1, failed_reason = $2 WHERE job_id = $3",
		types.JobStatusFailed, reason, jobID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s PsqlStore) GetFaceMatchByJobID(jobID string) (*types.JobRecord, error) {
	var faceMatchData types.JobRecord
	var completedAt, processedAt, failedAt sql.NullTime
	var failedReason sql.NullString
	err := s.db.QueryRow(
		"SELECT id, client_id, created_at, job_id, status, completed_at, processed_at, failed_at, failed_reason, match_score FROM face_match WHERE job_id = $1",
		jobID,
	).Scan(
		&faceMatchData.ID,
		&faceMatchData.ClientID,
		&faceMatchData.CreatedAt,
		&faceMatchData.JobID,
		&faceMatchData.Status,
		&completedAt,
		&processedAt,
		&failedAt,
		&failedReason,
		&faceMatchData.MatchScore,
	)
	if err != nil {
		return nil, err
	}

	// parsing the values
	faceMatchData.CompletedAt = parseTimeValue(completedAt)
	faceMatchData.ProcessedAt = parseTimeValue(processedAt)
	faceMatchData.FailedAt = parseTimeValue(failedAt)
	faceMatchData.FailedReason = parseStringValue(failedReason)

	// setting the type of job
	faceMatchData.Type = types.FaceMatchWorkType
	return &faceMatchData, nil
}

func (s PsqlStore) GetOCRByJobID(jobID string) (*types.JobRecord, error) {
	var ocrData types.JobRecord
	var completedAt, processedAt, failedAt sql.NullTime
	var failedReason sql.NullString
	var rawOCRDetails json.RawMessage
	err := s.db.QueryRow(
		"SELECT id, client_id, created_at, job_id, status, completed_at, processed_at, failed_at, failed_reason, details FROM ocr WHERE job_id = $1",
		jobID,
	).Scan(
		&ocrData.ID,
		&ocrData.ClientID,
		&ocrData.CreatedAt,
		&ocrData.JobID,
		&ocrData.Status,
		&completedAt,
		&processedAt,
		&failedAt,
		&failedReason,
		&rawOCRDetails,
	)
	if err != nil {
		return nil, err
	}

	// setting the type of job
	ocrData.Type = types.OCRWorkType

	// parsing the values
	ocrData.CompletedAt = parseTimeValue(completedAt)
	ocrData.ProcessedAt = parseTimeValue(processedAt)
	ocrData.FailedAt = parseTimeValue(failedAt)
	ocrData.FailedReason = parseStringValue(failedReason)

	// need to unmarshal the raw details saved in jsonb
	if err := json.Unmarshal(rawOCRDetails, &ocrData.OCRDetails); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ocr details: %v", err)
	}

	return &ocrData, nil
}

func parseTimeValue(dbTime sql.NullTime) string {
	if dbTime.Valid {
		return dbTime.Time.Format("2006-01-02 15:04:05.000000")
	} else {
		return "NULL"
	}
}

func parseStringValue(dnString sql.NullString) string {
	if dnString.Valid {
		return dnString.String
	} else {
		return "NULL"
	}
}
