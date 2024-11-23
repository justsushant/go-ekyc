package worker

import (
	"time"

	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
)

type OCRService struct{}

func NewOCRService() *OCRService {
	return &OCRService{}
}

func (d *OCRService) PerformOCR(payload types.OCRPayload) (*types.OCRResponse, error) {
	time.Sleep(5 * time.Second)
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
