package types

type ClientData struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	PlanID        int    `json:"plan_id"`
	AccessKey     string `json:"access_key"`
	SecretKeyHash string `json:"secret_key_hash"`
}

type UploadMetaData struct {
	Type       string `json:"type"`
	ClientID   int    `json:"client_id"`
	FilePath   string `json:"file_path"`
	FileSizeKB int64  `json:"file_size_kb"`
}

type FaceMatchPayload struct {
	ImageID1 string `json:"image1"`
	ImageID2 string `json:"image2"`
}

type OCRPayload struct {
	ImageID string `json:"image"`
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
