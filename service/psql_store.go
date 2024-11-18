package service

import (
	"database/sql"

	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/db"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
)

type PsqlStore struct {
	db *sql.DB
}

func NewPsqlStore(dsn string) PsqlStore {
	psqlClient := db.NewPostgreSQLStorage(dsn)
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
		"SELECT type, client_id, file_path, file_size_kb FROM client WHERE file_path = $1",
		imgUuid,
	).Scan(
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
