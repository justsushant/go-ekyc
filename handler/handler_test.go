package handler

import (
	"bytes"
	"encoding/json"
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

func (m mockService) ValidateImage(payload types.FaceMatchPayload) error {
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
