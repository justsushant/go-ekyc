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
