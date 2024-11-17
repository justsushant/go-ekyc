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

func (s PsqlStore) InsertClientData(planId int, payload types.SignupPayload) error {
	_, err := s.db.Exec("INSERT INTO client (name, email, plan_id) VALUES ($1, $2, $3)", payload.Name, payload.Email, planId)
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

func (s PsqlStore) UpdateClientKey(clientId int, accessKey, secretKeyHash string) error {
	_, err := s.db.Exec("UPDATE client access_key = '$1', secret_key_hash = '@2' WHERE id = '$3'", accessKey, secretKeyHash, clientId)
	if err != nil {
		return err
	}
	return nil
}
