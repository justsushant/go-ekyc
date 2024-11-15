package store

import "github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"

type DataStore interface {
	GetPlanIdFromName(planName string) (int, error)
	InsertClientData(payload types.SignupPayload, planId int, refreshToken string) error
}
