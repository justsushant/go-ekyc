package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"testing"

	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
)

var (
	Host                   = "http://localhost:8080"
	AccessKey              string
	SecretKey              string
	OCRFileUploadID        string
	OCRJobID               string
	FaceMatchFileUploadID1 string
	FaceMatchFileUploadID2 string
	FaceMatchJobID         string
)

// TODO: Setup proper integration test in docker containers which can be destri=oyed upon usage
func TestIntegrationMain(t *testing.T) {
	t.Run("Test health check endpoint", func(t *testing.T) {
		// arrange
		var expResp = types.HealthResponse{
			Message: "OK",
		}

		// make http request
		url := fmt.Sprintf("%s%s", Host, "/api/v1/health")
		resp := makeGetRequest(t, url, nil)
		defer resp.Body.Close()

		// check status code
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected %d but got %d", http.StatusOK, resp.StatusCode)
		}

		// check json body
		var respBody types.HealthResponse
		err := json.NewDecoder(resp.Body).Decode(&respBody)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if !reflect.DeepEqual(respBody, expResp) {
			t.Errorf("Expected %v but got %v", expResp, respBody)
		}
	})

	t.Run("Test signup endpoint", func(t *testing.T) {
		// arrange
		payload := types.SignupPayload{
			Name:  "max",
			Email: "max@one2n.in",
			Plan:  "basic",
		}
		body, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("Error while marshalling payload: %v", err)
		}

		// act
		url := fmt.Sprintf("%s%s", Host, "/api/v1/signup")
		resp := makePostRequest(t, url, nil, bytes.NewReader(body))
		defer resp.Body.Close()

		// check status code
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected %d but got %d", http.StatusOK, resp.StatusCode)
		}

		// check json body
		var respBody types.SignupResponse
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		if respBody.AccessKey == "" {
			t.Errorf("Access key empty: %v", respBody.AccessKey)
		}
		if respBody.SecretKey == "" {
			t.Errorf("Secret key empty: %v", respBody.AccessKey)
		}

		// saving accessKey and secretKey for usage in subsequent requests
		AccessKey = respBody.AccessKey
		SecretKey = respBody.SecretKey
	})

	// simulates the following scenario:
	// - upload image
	// - creates the ocr async operation
	// - checks the result
	testOCRIntegration(t)

	// simulates the following scenario:
	// - uploads two images
	// - creates the face async operation
	// - checks the result
	testFaceMatchIntegration(t)

	// tests the cache
	// checks if we get same jobID for same upload ids
	testCache(t)

}

func testOCRIntegration(t *testing.T) {
	t.Run("Test OCR file upload endpoint", func(t *testing.T) {
		// arrange
		url := fmt.Sprintf("%s%s", Host, "/api/v1/upload")
		textFields := map[string]string{"type": "id_card"}
		fileFields := map[string]string{"file": "../testdata/sample_image_1.jpeg"}
		body, contentType := makeMultiFormReqBody(t, textFields, fileFields)
		headers := map[string]string{"Content-Type": contentType, "accessKey": AccessKey, "secretKey": SecretKey}

		// make http request
		resp := makePostRequest(t, url, headers, body)
		defer resp.Body.Close()

		// check status code
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected %d but got %d", http.StatusOK, resp.StatusCode)
		}

		// check json body
		var respBody types.FileUploadResponse
		err := json.NewDecoder(resp.Body).Decode(&respBody)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		if respBody.Id == "" {
			t.Errorf("File upload id empty: %v", respBody.Id)
		}

		// saving ocr file upload id to be used in subsequent requests
		OCRFileUploadID = respBody.Id
	})

	t.Run("Test OCR Async operation", func(t *testing.T) {
		// arrange
		payload := types.OCRPayload{
			Image: OCRFileUploadID,
		}
		requestBody, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("Error while marshalling payload: %v", err)
		}

		// act
		url := fmt.Sprintf("%s%s", Host, "/api/v1/ocr-async")
		headers := map[string]string{"Content-Type": "application/json", "accessKey": AccessKey, "secretKey": SecretKey}
		resp := makePostRequest(t, url, headers, bytes.NewReader(requestBody))
		defer resp.Body.Close()

		// check status code
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected %d but got %d", http.StatusOK, resp.StatusCode)
		}

		// check json body
		var respBody types.OCRAsyncResponse
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		if respBody.Id == "" {
			t.Errorf("Job Id empty: %v", respBody.Id)
		}

		// saving jobID for usage in subsequent requests
		OCRJobID = respBody.Id
	})

	t.Run("Test Result for OCR Async operation", func(t *testing.T) {
		// act
		url := fmt.Sprintf("%s%s%s", Host, "/api/v1/result/ocr/", OCRJobID)
		header := map[string]string{"Content-Type": "application/json", "accessKey": AccessKey, "secretKey": SecretKey}
		resp := makeGetRequest(t, url, header)
		defer resp.Body.Close()

		// check status code
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected %d but got %d", http.StatusOK, resp.StatusCode)
		}

		// check json body
		var respBody types.ResultResponse
		err := json.NewDecoder(resp.Body).Decode(&respBody)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		if respBody.Status == "" {
			t.Errorf("Status empty: %v", respBody.Status)
		}
	})
}

func testFaceMatchIntegration(t *testing.T) {
	t.Run("Test face match file upload endpoint1", func(t *testing.T) {
		// arrange
		url := fmt.Sprintf("%s%s", Host, "/api/v1/upload")
		textFields := map[string]string{"type": "face"}
		fileFields := map[string]string{"file": "../testdata/sample_image_1.jpeg"}
		body, contentType := makeMultiFormReqBody(t, textFields, fileFields)
		headers := map[string]string{"Content-Type": contentType, "accessKey": AccessKey, "secretKey": SecretKey}

		// make http request
		resp := makePostRequest(t, url, headers, body)
		defer resp.Body.Close()

		// check status code
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected %d but got %d", http.StatusOK, resp.StatusCode)
		}

		// check json body
		var respBody types.FileUploadResponse
		err := json.NewDecoder(resp.Body).Decode(&respBody)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		if respBody.Id == "" {
			t.Errorf("File upload id empty: %v", respBody.Id)
		}

		// saving ocr file upload id to be used in subsequent requests
		FaceMatchFileUploadID1 = respBody.Id
	})

	t.Run("Test face match file upload endpoint2", func(t *testing.T) {
		// arrange
		url := fmt.Sprintf("%s%s", Host, "/api/v1/upload")
		textFields := map[string]string{"type": "face"}
		fileFields := map[string]string{"file": "../testdata/sample_image_1.jpeg"}
		body, contentType := makeMultiFormReqBody(t, textFields, fileFields)
		headers := map[string]string{"Content-Type": contentType, "accessKey": AccessKey, "secretKey": SecretKey}

		// make http request
		resp := makePostRequest(t, url, headers, body)
		defer resp.Body.Close()

		// check status code
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected %d but got %d", http.StatusOK, resp.StatusCode)
		}

		// check json body
		var respBody types.FileUploadResponse
		err := json.NewDecoder(resp.Body).Decode(&respBody)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		if respBody.Id == "" {
			t.Errorf("File upload id empty: %v", respBody.Id)
		}

		// saving ocr file upload id to be used in subsequent requests
		FaceMatchFileUploadID2 = respBody.Id
	})

	t.Run("Test Face Match Async operation", func(t *testing.T) {
		// arrange
		payload := types.FaceMatchPayload{
			Image1: FaceMatchFileUploadID1,
			Image2: FaceMatchFileUploadID2,
		}
		requestBody, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("Error while marshalling payload: %v", err)
		}

		// act
		url := fmt.Sprintf("%s%s", Host, "/api/v1/face-match-async")
		headers := map[string]string{"Content-Type": "application/json", "accessKey": AccessKey, "secretKey": SecretKey}
		resp := makePostRequest(t, url, headers, bytes.NewReader(requestBody))
		defer resp.Body.Close()

		// check status code
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected %d but got %d", http.StatusOK, resp.StatusCode)
		}

		// check json body
		var respBody types.FaceMatchAsyncResponse
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		if respBody.Id == "" {
			t.Errorf("Job Id empty: %v", respBody.Id)
		}

		// saving jobID for usage in subsequent requests
		FaceMatchJobID = respBody.Id
	})

	t.Run("Test Result for Face Match Async operation", func(t *testing.T) {
		// act
		url := fmt.Sprintf("%s%s%s", Host, "/api/v1/result/face_match/", FaceMatchJobID)
		header := map[string]string{"Content-Type": "application/json", "accessKey": AccessKey, "secretKey": SecretKey}
		resp := makeGetRequest(t, url, header)
		defer resp.Body.Close()

		// check status code
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected %d but got %d", http.StatusOK, resp.StatusCode)
		}

		// check json body
		var respBody types.ResultResponse
		err := json.NewDecoder(resp.Body).Decode(&respBody)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		if respBody.Status == "" {
			t.Errorf("Status empty: %v", respBody.Status)
		}
	})
}

func testCache(t *testing.T) {
	t.Run("Test the cache for OCR", func(t *testing.T) {
		// arrange
		payload := types.OCRPayload{
			Image: OCRFileUploadID,
		}
		requestBody, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("Error while marshalling payload: %v", err)
		}

		// act
		url := fmt.Sprintf("%s%s", Host, "/api/v1/ocr-async")
		headers := map[string]string{"Content-Type": "application/json", "accessKey": AccessKey, "secretKey": SecretKey}
		resp := makePostRequest(t, url, headers, bytes.NewReader(requestBody))
		defer resp.Body.Close()

		// check json body
		var respBody types.OCRAsyncResponse
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		if respBody.Id != OCRJobID {
			t.Errorf("Expected the jobID to be %q but got %q", OCRJobID, respBody.Id)
		}
	})

	t.Run("Test the cache for face match", func(t *testing.T) {
		// arrange
		payload := types.FaceMatchPayload{
			Image1: FaceMatchFileUploadID1,
			Image2: FaceMatchFileUploadID2,
		}
		requestBody, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("Error while marshalling payload: %v", err)
		}

		// act
		url := fmt.Sprintf("%s%s", Host, "/api/v1/face-match-async")
		headers := map[string]string{"Content-Type": "application/json", "accessKey": AccessKey, "secretKey": SecretKey}
		resp := makePostRequest(t, url, headers, bytes.NewReader(requestBody))
		defer resp.Body.Close()

		// check json body
		var respBody types.FaceMatchAsyncResponse
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		if respBody.Id != FaceMatchJobID {
			t.Errorf("Expected the jobID to be %q but got %q", FaceMatchJobID, respBody.Id)
		}
	})
}

func makeGetRequest(t *testing.T, url string, headers map[string]string) *http.Response {
	t.Helper()

	// make request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("Error while making GET request on %s: %s", url, err.Error())
		return nil
	}

	// setting headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Error while sending GET request on %s: %s", url, err.Error())
		return nil
	}

	return resp
}

func makePostRequest(t *testing.T, url string, headers map[string]string, body io.Reader) *http.Response {
	// make request
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		t.Fatalf("Error while making POST request on %s: %s", url, err.Error())
		return nil
	}

	// setting headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Error while sending GET request on %s: %s", url, err.Error())
		return nil
	}

	return resp
}

func makeMultiFormReqBody(t *testing.T, fields map[string]string, data map[string]string) (io.Reader, string) {
	t.Helper()

	// to hold the multipart form data
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// writing text fields
	for k, v := range fields {
		err := writer.WriteField(k, v)
		if err != nil {
			t.Fatalf("Error while writing fields to multi part form data: %s", err.Error())
		}
	}

	// file path: value, key name: key
	// writing file fields
	for k, v := range data {
		// accessing file
		file, err := os.Open(v)
		if err != nil {
			t.Fatalf("Error while opening file on the path %s: %s", v, err.Error())
		}
		defer file.Close()

		// creating form field
		part, err := writer.CreateFormFile(k, "image.jpeg") // "file" is the key in the form
		if err != nil {
			t.Fatalf("Error while creating form file: %s", err.Error())
		}

		// copying file into form field
		_, err = io.Copy(part, file)
		if err != nil {
			t.Fatalf("Error while copying file in form field: %s", err.Error())
		}
	}
	contentType := writer.FormDataContentType()
	defer writer.Close()

	return &body, contentType
}
