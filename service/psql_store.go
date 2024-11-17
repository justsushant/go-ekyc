package service

import (
	"database/sql"

	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
)

type PsqlStore struct {
	db *sql.DB
}

func NewPsqlStore(db *sql.DB) PsqlStore {
	return PsqlStore{
		db: db,
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
