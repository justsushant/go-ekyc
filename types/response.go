package types

import (
	"encoding/json"
	"log"
)

type ErrorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

type HealthResponse struct {
	Message string `json:"message"`
}

type SignupResponse struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}

type IDResponse struct {
	Id string `name:"id"`
}

type FileUploadResponse IDResponse
type OCRAsyncResponse IDResponse
type FaceMatchAsyncResponse IDResponse

type ResultResponse struct {
	Status       string      `json:"status"`
	Message      string      `json:"message"`
	CreatedAt    string      `json:"created_at"`
	ProcessedAt  string      `json:"processed_at"`
	CompletedAt  string      `json:"completed_at"`
	FailedAt     string      `json:"failed_at"`
	FailedReason string      `json:"failed_reason"`
	Result       OCRResponse `json:"result"`
}

type FaceMatchResponse int

type OCRResponseRaw json.RawMessage

type OCRResponse struct {
	Name      string `json:"name"`
	Gender    string `json:"gender"`
	DOB       string `json:"dateOfBirth"`
	IdNumber  string `json:"idNumber"`
	AddrLine1 string `json:"addressLine1"`
	AddrLine2 string `json:"addressLine2"`
	Pincode   string `json:"pincode"`
}

func (or *OCRResponse) String() string {
	jsonData, err := json.Marshal(or)
	if err != nil {
		log.Fatal("Error while marshalling ocr response: ", err)
	}

	return string(jsonData)
}
