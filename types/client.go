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

type JobRecord struct {
	Type          WorkType          `json:"job_type"`
	ID            int               `json:"id"`
	ClientID      int               `json:"client_id"`
	CreatedAt     string            `json:"created_at"`
	JobID         string            `json:"job_id"`
	Status        string            `json:"status"`
	CompletedAt   string            `json:"completed_at"`
	ProcessedAt   string            `json:"processed_at"`
	FailedAt      string            `json:"failed_at"`
	FailedReason  string            `json:"failed_reason"`
	MatchScore    FaceMatchResponse `json:"match_score"`
	RawOCRDetails json.RawMessage   `json:"details"`
	OCRDetails    OCRResponse       `json:"-"`
}

type ClientReport struct {
	ClientID          string `csv:"client_id"`
	Name              string `csv:"name"`
	Plan              string `csv:"plan"`
	Date              string `csv:"date"`
	TotalFaceMatch    string `csv:"total_face_match_for_day"`
	TotalOcr          string `csv:"total_ocr_for_da"`
	TotalImgStorageMB string `csv:"total_image_storage_in_mb"`
	TotalAPIUsageCost string `csv:"api_usage_cost_usd"`
	TotalStorageCost  string `csv:"storage_cost_usd"`
}

// type ClientReport struct {
// 	ClientID          string `json:"client_id"`
// 	Name              string `json:"name"`
// 	Plan              string `json:"plan"`
// 	Date              string `json:"date"`
// 	TotalFaceMatch    string `json:"total_face_match_for_day"`
// 	TotalOcr          string `json:"total_ocr_for_da"`
// 	TotalImgStorageMB string `json:"total_image_storage_in_mb"`
// 	TotalAPIUsageCost string `json:"api_usage_cost_usd"`
// 	TotalStorageCost  string `json:"storage_cost_usd"`
// }
