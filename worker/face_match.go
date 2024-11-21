package worker

import (
	"math/rand"

	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
)

type FaceMatchService struct{}

func NewFaceMatchService() *FaceMatchService {
	return &FaceMatchService{}
}

func (d *FaceMatchService) CalcFaceMatchScore(payload types.FaceMatchPayload) (int, error) {
	// time.Sleep(5 * time.Second)
	return rand.Intn(100) + 1, nil
}
