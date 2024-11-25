package types

type Plan string
type WorkType string
type JobStatus string
type FileType string

const ()

const (
	BASIC_PLAN      = "basic"
	ADVANCE_PLAN    = "advance"
	ENTERPRISE_PLAN = "enterprise"

	FACE_MATCH_WORK_TYPE = "face_match"
	OCR_WORK_TYPE        = "ocr"

	JOB_STATUS_PROCESSING = "processing"
	JOB_STATUS_CREATED    = "created"
	JOB_STATUS_COMPLETED  = "completed"
	JOB_STATUS_FAILED     = "failed"

	FACE_TYPE    = "face"
	ID_CARD_TYPE = "id_card"

	VALID_FORMAT_PNG  = ".png"
	VALID_FORMAT_JPEG = ".jpeg"
)
