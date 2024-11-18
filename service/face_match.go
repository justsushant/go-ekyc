package service

import (
	"math/rand"

	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
)

type FaceMatcher interface {
	CalcFaceMatchScore(payload types.FaceMatchPayload) (int, error)
}

type DummyFaceMatchService struct{}

func (d *DummyFaceMatchService) CalcFaceMatchScore(payload types.FaceMatchPayload) (int, error) {
	return rand.Intn(100) + 1, nil
}
