package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/service"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
	"github.com/stretchr/testify/assert"
)

// mock of client service for usage in tests
type mockService struct{}

func (m mockService) ValidatePayload(payload types.SignupPayload) error {
	if payload.Email == "test@abc@corp" {
		return service.ErrInvalidEmail
	} else if payload.Plan == "invalid-plan" {
		return service.ErrInvalidPlan
	} else {
		return nil
	}
}

func (m mockService) GenerateKeyPair() (*service.KeyPair, error) {
	return service.NewKeyPair("qwerty", "quirkyfox", ""), nil
}

func (m mockService) SaveSignupData(payload types.SignupPayload, keyPair *service.KeyPair) error {
	return nil
}

func (m mockService) ValidateFile(fileName, fileType string) error {
	if fileType != "face" && fileType != "id_card" {
		return service.ErrInvalidFileType
	}

	ext := filepath.Ext(fileName)
	if ext != types.VALID_FORMAT_PNG && ext != types.VALID_FORMAT_JPEG {
		return service.ErrInvalidFileFormat
	}

	return nil
}

func (m mockService) SaveFile(fileHeader *multipart.FileHeader, uploadMetaData *types.UploadMetaData) error {
	return nil
}

func (m mockService) ValidateImage(payload types.FaceMatchPayload, clientID int) error {
	if payload.Image1 == "exec" {
		return service.ErrInvalidImgId
	}

	if payload.Image2 == "qwerty" {
		return service.ErrNotFaceImg
	}

	return nil
}

func (m mockService) CalcAndSaveFaceMatchScore(payload types.FaceMatchPayload, clientID int) (int, error) {
	return 72, nil
}

func (m mockService) ValidateImageOCR(payload types.OCRPayload, clientID int) error {
	if payload.Image == "invalid-img" {
		return service.ErrInvalidImgId
	}

	if payload.Image == "not-id-card" {
		return service.ErrNotIDCardImg
	}

	return nil
}

func (m mockService) PerformAndSaveOCR(payload types.OCRPayload, clientID int) (*types.OCRResponse, error) {
	return &types.OCRResponse{
		Name:      "John Adams",
		Gender:    "Male",
		DOB:       "1990-01-24",
		IdNumber:  "1234-1234-1234",
		AddrLine1: "A2, 201, Amar Villa",
		AddrLine2: "MG Road, Pune",
		Pincode:   "411004",
	}, nil
}

func (m mockService) PerformFaceMatchAsync(payload types.FaceMatchPayload, clientID int) (string, error) {
	if payload.Image1 == "exec" {
		return "", service.ErrInvalidImgId
	}

	if payload.Image2 == "qwerty" {
		return "", service.ErrNotFaceImg
	}

	return "uuid-ok", nil
}

func (m mockService) PerformOCRAsync(payload types.OCRPayload, clientID int) (string, error) {
	if payload.Image == "invalid-img" {
		return "", service.ErrInvalidImgId
	}
	if payload.Image == "not-id-card" {
		return "", service.ErrNotIDCardImg
	}

	return "uuid-ok", nil
}
func (m mockService) FetchDataFromCache(payload interface{}, clientID int, jobType string) (string, bool) {
	return "", false
}
func (m mockService) SetDataInCache(payload interface{}, clientID int, jobType, jobID string) {}

func (m mockService) GetJobDetailsByJobID(jobID, jobType string) (*types.JobRecord, error) {
	if jobType == "face_match" && jobID == "qwerty" {
		return &types.JobRecord{
			ClientID: 2,
			Type:     types.FaceMatchWorkType,
		}, nil
	}

	if jobType == "ocr" && jobID == "zxcvb" {
		return &types.JobRecord{
			ClientID:    1,
			Status:      types.JobStatusProcessing,
			ProcessedAt: "timestamp",
		}, nil
	}

	if jobType == "ocr" && jobID == "zxcvba" {
		return &types.JobRecord{
			ClientID:  1,
			Status:    types.JobStatusCreated,
			CreatedAt: "timestamp",
			Type:      types.OCRWorkType,
		}, nil
	}

	if jobType == "ocr" && jobID == "zgcvba" {
		return &types.JobRecord{
			ClientID:     1,
			Status:       types.JobStatusFailed,
			FailedAt:     "timestamp",
			FailedReason: "reason",
			Type:         types.OCRWorkType,
		}, nil
	}

	if jobType == "face_match" && jobID == "mnlkpo" {
		return &types.JobRecord{
			ClientID:    1,
			Status:      types.JobStatusCompleted,
			CompletedAt: "timestamp",
			MatchScore:  72,
			Type:        types.FaceMatchWorkType,
		}, nil
	}

	if jobType == "ocr" && jobID == "mnlkpo" {
		return &types.JobRecord{
			ClientID:    1,
			Status:      types.JobStatusCompleted,
			CompletedAt: "timestamp",
			Type:        types.OCRWorkType,
			OCRDetails: types.OCRResponse{
				Name:      "xyz",
				Gender:    "xyz",
				DOB:       "xyz",
				IdNumber:  "xyz",
				AddrLine1: "xyz",
				AddrLine2: "xyz",
				Pincode:   "xyz",
			},
		}, nil
	}

	return nil, nil
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
		{
			name: "valid case",
			payload: types.SignupPayload{
				Name:  "abc corp",
				Email: "test@abc.corp",
				Plan:  "basic",
			},
			expStatusCode: http.StatusOK,
			expResponse:   `{"accessKey":"qwerty", "secretKey":"quirkyfox"}`,
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
			handler := NewHandler(&mockService{})
			handler.SignupHandler(c)

			// asserting the values
			assert.Equal(t, tc.expStatusCode, w.Code)
			assert.JSONEq(t, tc.expResponse, w.Body.String())
		})
	}
}

func TestFileUploadHandler(t *testing.T) {
	tt := []struct {
		name          string
		fileName      string
		fileType      string
		content       string
		expStatusCode int
		expResponse   string
	}{
		{
			name:          "invalid file type case",
			fileName:      "invalid.jpeg",
			fileType:      "invalid-type",
			content:       "Hello, world!",
			expStatusCode: http.StatusBadRequest,
			expResponse:   `{"errorMessage": "invalid type, supported types are face or id_card"}`,
		},
		{
			name:          "invalid empty file type case",
			fileName:      "invalid.jpeg",
			fileType:      "",
			content:       "Hello, world!",
			expStatusCode: http.StatusBadRequest,
			expResponse:   `{"errorMessage": "invalid type, supported types are face or id_card"}`,
		},
		{
			name:          "invalid file name without ext case",
			fileName:      "fileName",
			fileType:      "face",
			content:       "Hello, world!",
			expStatusCode: http.StatusBadRequest,
			expResponse:   `{"errorMessage": "invalid file format, supported formats are png or jpeg"}`,
		},
		{
			name:          "invalid file format case",
			fileName:      "invalid.cyan",
			fileType:      "face",
			content:       "Hello, world!",
			expStatusCode: http.StatusBadRequest,
			expResponse:   `{"errorMessage": "invalid file format, supported formats are png or jpeg"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// reading file from the request body
			var buf bytes.Buffer

			writer := multipart.NewWriter(&buf)
			part, err := writer.CreateFormFile("file", tc.fileName)
			if err != nil {
				t.Fatalf("Error creating form file: %v", err)
			}
			part.Write([]byte(tc.content))

			// reading normal key-value pair
			writer.WriteField("type", tc.fileType)

			// closing the writer
			writer.Close()

			// preparing the test
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/image", &buf)
			c.Request.Header.Set("Content-Type", writer.FormDataContentType())

			// calling the file upload handler
			handler := NewHandler(&mockService{})
			handler.FileUploadHandler(c)

			// asserting the values
			assert.Equal(t, tc.expStatusCode, w.Code)
			assert.JSONEq(t, tc.expResponse, w.Body.String())
		})
	}
}

func TestFaceMatchHandler(t *testing.T) {
	tt := []struct {
		name          string
		payload       types.FaceMatchPayload
		expStatusCode int
		expResponse   string
	}{
		{
			name: "invalid img id case",
			payload: types.FaceMatchPayload{
				Image1: "exec",
				Image2: "qwerty-valid",
			},
			expStatusCode: http.StatusBadRequest,
			expResponse:   `{"errorMessage":"invalid or missing image id"}`,
		},
		{
			name: "invalid img type case",
			payload: types.FaceMatchPayload{
				Image1: "exec-valid",
				Image2: "qwerty",
			},
			expStatusCode: http.StatusBadRequest,
			expResponse:   `{"errorMessage":"not a face image"}`,
		},
		{
			name: "valid face match case",
			payload: types.FaceMatchPayload{
				Image1: "exec-valid",
				Image2: "qwerty-valid",
			},
			expStatusCode: http.StatusOK,
			expResponse:   `{"score":72}`,
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
			c.Request = httptest.NewRequest("POST", "/face-match", bytes.NewBuffer([]byte(body)))
			c.Request.Header.Set("Content-Type", "application/json")
			c.Set("client_id", 4)

			// calling the signup handler
			handler := NewHandler(&mockService{})
			handler.FaceMatchHandler(c)

			// asserting the values
			assert.Equal(t, tc.expStatusCode, w.Code)
			assert.JSONEq(t, tc.expResponse, w.Body.String())
		})
	}
}

func TestOCRHandler(t *testing.T) {
	tt := []struct {
		name          string
		payload       types.OCRPayload
		expStatusCode int
		expResponse   string
	}{
		{
			name: "invalid img id case",
			payload: types.OCRPayload{
				Image: "invalid-img",
			},
			expStatusCode: http.StatusBadRequest,
			expResponse:   `{"errorMessage":"invalid or missing image id"}`,
		},
		{
			name: "invalid img type case",
			payload: types.OCRPayload{
				Image: "not-id-card",
			},
			expStatusCode: http.StatusBadRequest,
			expResponse:   `{"errorMessage":"not an id card image"}`,
		},
		{
			name: "valid face match case",
			payload: types.OCRPayload{
				Image: "ac",
			},
			expStatusCode: http.StatusOK,
			expResponse: `{
				"name": "John Adams",
				"gender": "Male",
				"dateOfBirth": "1990-01-24",
				"idNumber": "1234-1234-1234",
				"addressLine1": "A2, 201, Amar Villa",
				"addressLine2": "MG Road, Pune",
				"pincode": "411004"
			}`,
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
			c.Request = httptest.NewRequest("POST", "/ocr", bytes.NewBuffer([]byte(body)))
			c.Request.Header.Set("Content-Type", "application/json")
			c.Set("client_id", 4)

			// calling the signup handler
			handler := NewHandler(&mockService{})
			handler.OCRHandler(c)

			// asserting the values
			assert.Equal(t, tc.expStatusCode, w.Code)
			assert.JSONEq(t, tc.expResponse, w.Body.String())
		})
	}
}

func TestFaceMatchHandlerAsync(t *testing.T) {
	tt := []struct {
		name          string
		payload       types.FaceMatchPayload
		expStatusCode int
		expResponse   string
	}{
		{
			name: "invalid img id case",
			payload: types.FaceMatchPayload{
				Image1: "exec",
				Image2: "qwerty-valid",
			},
			expStatusCode: http.StatusBadRequest,
			expResponse:   `{"errorMessage":"invalid or missing image id"}`,
		},
		{
			name: "invalid img type case",
			payload: types.FaceMatchPayload{
				Image1: "exec-valid",
				Image2: "qwerty",
			},
			expStatusCode: http.StatusBadRequest,
			expResponse:   `{"errorMessage":"not a face image"}`,
		},
		{
			name: "valid face match case",
			payload: types.FaceMatchPayload{
				Image1: "exec-valid",
				Image2: "qwerty-valid",
			},
			expStatusCode: http.StatusOK,
			expResponse:   `{"id":"uuid-ok"}`,
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
			c.Request = httptest.NewRequest("POST", "/face-match-async", bytes.NewBuffer([]byte(body)))
			c.Request.Header.Set("Content-Type", "application/json")
			c.Set("client_id", 4)

			// calling the signup handler
			handler := NewHandler(&mockService{})
			handler.FaceMatchHandlerAsync(c)

			// asserting the values
			assert.Equal(t, tc.expStatusCode, w.Code)
			assert.JSONEq(t, tc.expResponse, w.Body.String())
		})
	}
}

func TestOCRHandlerAsync(t *testing.T) {
	tt := []struct {
		name          string
		payload       types.OCRPayload
		expStatusCode int
		expResponse   string
	}{
		{
			name: "invalid img id case",
			payload: types.OCRPayload{
				Image: "invalid-img",
			},
			expStatusCode: http.StatusBadRequest,
			expResponse:   `{"errorMessage":"invalid or missing image id"}`,
		},
		{
			name: "invalid img type case",
			payload: types.OCRPayload{
				Image: "not-id-card",
			},
			expStatusCode: http.StatusBadRequest,
			expResponse:   `{"errorMessage":"not an id card image"}`,
		},
		{
			name: "valid face match case",
			payload: types.OCRPayload{
				Image: "ac",
			},
			expStatusCode: http.StatusOK,
			expResponse:   `{"id": "uuid-ok"}`,
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
			c.Request = httptest.NewRequest("POST", "/ocr-async", bytes.NewBuffer([]byte(body)))
			c.Request.Header.Set("Content-Type", "application/json")
			c.Set("client_id", 4)

			// calling the signup handler
			handler := NewHandler(&mockService{})
			handler.OCRHandlerAsync(c)

			// asserting the values
			assert.Equal(t, tc.expStatusCode, w.Code)
			assert.JSONEq(t, tc.expResponse, w.Body.String())
		})
	}
}

func TestResultHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tt := []struct {
		name          string
		jobType       string
		jobID         string
		clientID      int
		expStatusCode int
		expResponse   string
	}{
		{
			name:          "jobID for different client",
			jobType:       "face_match",
			jobID:         "qwerty",
			clientID:      1,
			expStatusCode: 400,
			expResponse:   `{ "errorMessage": "invalid or missing job id" }`,
		},
		{
			name:          "job is still processing",
			jobType:       "ocr",
			jobID:         "zxcvb",
			clientID:      1,
			expStatusCode: 200,
			expResponse:   `{ "status": "processing", "processed_at": "timestamp", "message": "job is still running"}`,
		},
		{
			name:          "job is created",
			jobType:       "ocr",
			jobID:         "zxcvba",
			clientID:      1,
			expStatusCode: 200,
			expResponse:   `{ "status": "created", "created_at": "timestamp", "message": "job is created"}`,
		},
		{
			name:          "job is failed",
			jobType:       "ocr",
			jobID:         "zgcvba",
			clientID:      1,
			expStatusCode: 200,
			expResponse:   `{ "status": "failed", "failed_at": "timestamp", "message": "job is failed", "failed_reason": "reason"}`,
		},
		{
			name:          "job is completed face_match",
			jobType:       "face_match",
			jobID:         "mnlkpo",
			clientID:      1,
			expStatusCode: 200,
			expResponse:   `{ "status": "completed", "completed_at": "timestamp", "message": "job is completed", "result": 72}`,
		},
		{
			name:          "job is completed ocr",
			jobType:       "ocr",
			jobID:         "mnlkpo",
			clientID:      1,
			expStatusCode: 200,
			expResponse:   `{"completed_at": "timestamp", "message": "job is completed", "result": {"name":"xyz","gender":"xyz","dateOfBirth":"xyz","idNumber":"xyz","addressLine1":"xyz","addressLine2":"xyz","pincode":"xyz"}, "status": "completed"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// preparing the test
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", fmt.Sprintf("/result/%s/%s", tc.jobType, tc.jobID), nil)
			c.Request.Header.Set("Content-Type", "application/json")
			c.Set("client_id", tc.clientID)
			c.Params = []gin.Param{
				{Key: "jobType", Value: tc.jobType},
				{Key: "jobID", Value: tc.jobID},
			}

			// calling the signup handler
			handler := NewHandler(&mockService{})
			handler.ResultHandler(c)

			// asserting the values
			assert.Equal(t, tc.expStatusCode, w.Code)
			assert.JSONEq(t, tc.expResponse, w.Body.String())
		})
	}
}
