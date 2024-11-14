package types

type SignupPayload struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Plan  string `json:"plan"`
}
