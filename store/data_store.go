package store

import (
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
)

type DataStore interface {
	GetPlanIdFromName(planName string) (int, error)
	InsertClientData(planId int, payload types.SignupPayload, accessKey, secretKeyHash string) error
}
