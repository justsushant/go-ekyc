package types

import "io"

type Plan string

const BasicPlan = "basic"
const AdvancePlan = "advance"
const EnterprisePlan = "enterprise"

type WorkType string

const FaceMatchWorkType = "face_match"
const OCRWorkType = "ocr"

type JobStatus string

const JobStatusProcessing = "processing"
const JobStatusCreated = "created"
const JobStatusCompleted = "completed"
const JobStatusFailed = "failed"

type FileUpload struct {
	Name    string
	Content io.Reader
	Size    int64
	Headers map[string]string
}
