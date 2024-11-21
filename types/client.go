package types

import (
	"encoding/json"
	"log"
)

type ClientData struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	PlanID        int    `json:"plan_id"`
	AccessKey     string `json:"access_key"`
	SecretKeyHash string `json:"secret_key_hash"`
}

type UploadMetaData struct {
	Id         int    `json:"id"`
	Type       string `json:"type"`
	ClientID   int    `json:"client_id"`
	FilePath   string `json:"file_path"`
	FileSizeKB int64  `json:"file_size_kb"`
}

type FaceMatchPayload struct {
	Image1 string `json:"image1"`
	Image2 string `json:"image2"`
}

type FaceMatchData struct {
	ClientID int `json:"client_id"`
	ImageID1 int `json:"upload_id1"`
	ImageID2 int `json:"upload_id2"`
	Score    int `json:"score"`
}

type OCRPayload struct {
	Image string `json:"image"`
}

type OCRData struct {
	ClientID int    `json:"client_id"`
	ImageID  int    `json:"upload_id"`
	Data     string `json:"details"`
}

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

type OCRResult struct {
	ImageID  string      `json:"image"`
	ClientID int         `json:"client_id"`
	Data     OCRResponse `json:"details"`
}

type FaceMatchInternalPayload struct {
	JobID  string `json:"job_id"`
	Image1 string `json:"image1"`
	Image2 string `json:"image2"`
}

type FaceMatchQueuePayload struct {
	Type WorkType `json:"type"`
	Msg  FaceMatchInternalPayload
}

type OCRInternalPayload struct {
	JobID string `json:"job_id"`
	Image string `json:"image"`
}

type OCRQueuePayload struct {
	Type WorkType `json:"type"`
	Msg  OCRInternalPayload
}
