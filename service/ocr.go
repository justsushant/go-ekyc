package service

import (
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
)

type OCRPerformer interface {
	PerformOCR(payload types.OCRPayload) (*types.OCRResponse, error)
}

type DummyOcrService struct{}

func (d *DummyOcrService) PerformOCR(payload types.OCRPayload) (*types.OCRResponse, error) {
	return &types.OCRResponse{
		Name:      "John Adams",
		Gender:    "Male",
		DOB:       "1990-01-24",
		IdNumber:  "1234-1234-1234",
		AddrLine1: "A2, 201, Amar Villa",
		AddrLine2: "MG Road, Pune",
		Pincode:   "411004",
	}, nil
}
