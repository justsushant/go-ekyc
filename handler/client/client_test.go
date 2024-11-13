package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/controller/client"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
	"github.com/stretchr/testify/assert"
)

type mockClientService struct{}

func (m mockClientService) ValidatePayload(payload types.SignupPayload) error {
	if payload.Email == "test@abc@corp" {
		return client.ErrInvalidEmail
	} else if payload.Plan == "invalid-plan" {
		return client.ErrInvalidPlan
	} else {
		return nil
	}
}

func TestSignupHandler(t *testing.T) {
	tt := []struct {
		name          string
		payload       types.SignupPayload
		expStatusCode int
		expResponse   string
	}{
		{
			name: "invalid email case",
			payload: types.SignupPayload{
				Name:  "abc corp",
				Email: "test@abc@corp",
				Plan:  "basic",
			},
			expStatusCode: http.StatusBadRequest,
			expResponse:   `{"errorMessage":"invalid email"}`,
		},
		{
			name: "invalid plan case",
			payload: types.SignupPayload{
				Name:  "abc corp",
				Email: "test@abc.corp",
				Plan:  "invalid-plan",
			},
			expStatusCode: http.StatusBadRequest,
			expResponse:   `{"errorMessage":"invalid plan, supported plans are basic, advanced, or enterprise"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// marhalling the payload into json
			body, err := json.Marshal(tc.payload)
			if err != nil {
				t.Fatalf("Error while marshalling payload: %v", err)
			}

			// preparing the test
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/signup", bytes.NewBuffer([]byte(body)))
			c.Request.Header.Set("Content-Type", "application/json")

			// calling the signup handler
			handler := NewHandler(&mockClientService{})
			handler.SignupHandler(c)

			// asserting the values
			assert.Equal(t, tc.expStatusCode, w.Code)
			assert.JSONEq(t, tc.expResponse, w.Body.String())
		})
	}
}
