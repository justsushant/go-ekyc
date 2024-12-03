package types

type SignupPayload struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Plan  string `json:"plan"`
}

type FaceMatchPayload struct {
	Image1 string `json:"image1"`
	Image2 string `json:"image2"`
}

type OCRPayload struct {
	Image string `json:"image"`
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

type ResultPayload struct {
	ID string `json:"id"`
}
