package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
)

func TestSignupHandlerForInvalidEmail(t *testing.T) {
	// unhappy path of invalid email
	payload := types.SignupPayload{
		Name:  "abc corp",
		Email: "test@abc@corp",
		Plan:  "basic",
	}
	expStatus := http.StatusBadRequest
	// expResp := `{
	// 	"errorMessage”: "invalid email"
	// }`
	expErrorMessage := "invalid email"

	// marhalling the payload into json
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Error while marshalling payload: %v", err)
	}

	// preparing the test
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/signup", bytes.NewBuffer([]byte(body)))
	c.Request.Header.Set("Content-Type", "application/json")

	// calling the signup handler
	SignupHandler(c)

	// asserting the values
	if expStatus != w.Code {
		t.Errorf("Expected %d but got %d", expStatus, w.Code)
	}

	var respBody map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &respBody)
	if err != nil {
		t.Fatalf("Error while unmarshalling response: %v", err)
	}

	if expErrorMessage != respBody["errorMessage"] {
		t.Errorf("Expected %v but got %v", expStatus, w.Code)
	}

}
func TestSignupHandlerForInvalidPlan(t *testing.T) {
	// unhappy path of invalid email
	payload := types.SignupPayload{
		Name:  "abc corp",
		Email: "test@abc.corp",
		Plan:  "invalid-plan",
	}
	expStatus := http.StatusBadRequest
	expErrorMessage := "invalid plan, supported plans are basic, advanced, or enterprise"

	// marhalling the payload into json
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Error while marshalling payload: %v", err)
	}

	// preparing the test
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/signup", bytes.NewBuffer([]byte(body)))
	c.Request.Header.Set("Content-Type", "application/json")

	// calling the signup handler
	SignupHandler(c)

	// asserting the values
	if expStatus != w.Code {
		t.Errorf("Expected %d but got %d", expStatus, w.Code)
	}

	var respBody map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &respBody)
	if err != nil {
		t.Fatalf("Error while unmarshalling response: %v", err)
	}

	if expErrorMessage != respBody["errorMessage"] {
		t.Errorf("Expected %v but got %v", expStatus, w.Code)
	}

}
