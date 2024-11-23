package service

import (
	"errors"
	"reflect"
	"testing"

	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
	amqp "github.com/rabbitmq/amqp091-go"
)

type mockDataStore struct{}

func (m *mockDataStore) GetPlanIdFromName(planName string) (int, error) { return 0, nil }
func (m *mockDataStore) InsertClientData(planId int, payload types.SignupPayload, accessKey, secretKeyHash string) error {
	return nil
}

func (m *mockDataStore) GetClientFromAccessKey(accessKey string) (*types.ClientData, error) {
	return nil, nil
}

func (m *mockDataStore) InsertUploadMetaData(uploadMetaData *types.UploadMetaData) error { return nil }

func (m *mockDataStore) GetMetaDataByUUID(imgUuid string) (*types.UploadMetaData, error) {
	if imgUuid == "abc" {
		return &types.UploadMetaData{
			Type:     "not-face",
			ClientID: 1,
			FilePath: imgUuid,
		}, nil
	}
	if imgUuid == "xyz" {
		return &types.UploadMetaData{
			Type:     "face",
			ClientID: 1,
			FilePath: imgUuid,
		}, nil
	}
	if imgUuid == "pqr" {
		return &types.UploadMetaData{
			Type:     "not-face",
			ClientID: 2,
			FilePath: imgUuid,
		}, nil
	}
	if imgUuid == "def" {
		return &types.UploadMetaData{
			Type:     "face",
			ClientID: 2,
			FilePath: imgUuid,
		}, nil
	}
	if imgUuid == "ert" {
		return &types.UploadMetaData{}, nil
	}
	if imgUuid == "jnk" {
		return &types.UploadMetaData{}, nil
	}
	if imgUuid == "cvbas" {
		return &types.UploadMetaData{
			Type:     "not-id_card",
			ClientID: 2,
			FilePath: imgUuid,
		}, nil
	}
	if imgUuid == "cvbasrt" || imgUuid == "ac" {
		return &types.UploadMetaData{
			Type:     "id_card",
			ClientID: 3,
			FilePath: imgUuid,
		}, nil
	}
	if imgUuid == "asdf" || imgUuid == "lkjh" {
		return &types.UploadMetaData{
			Type:     "face",
			ClientID: 3,
			FilePath: imgUuid,
		}, nil
	}

	return nil, nil
}

func (m *mockDataStore) InsertFaceMatchResult(result *types.FaceMatchData) error { return nil }
func (m *mockDataStore) InsertOCRResult(result *types.OCRData) error             { return nil }
func (m *mockDataStore) InsertFaceMatchJobCreated(id1, id2, clientID int, jobID string) error {
	return nil
}
func (m *mockDataStore) InsertOCRJobCreated(id1, client_id int, jobID string) error {
	return nil
}

func (m *mockDataStore) UpdateFaceMatchJobCompleted(jobID string, score int) error {
	return nil
}

func (m *mockDataStore) UpdateOCRJobCompleted(jobID string, data *types.OCRResponse) error {
	return nil
}

func (m *mockDataStore) UpdateFaceMatchJobProcessed(jobID string) error             { return nil }
func (m *mockDataStore) UpdateOCRJobProcessed(jobID string) error                   { return nil }
func (m *mockDataStore) UpdateFaceMatchJobFailed(jobID, reason string) error        { return nil }
func (m *mockDataStore) UpdateOCRJobFailed(jobID, reason string) error              { return nil }
func (m *mockDataStore) GetFaceMatchByJobID(jobID string) (*types.JobRecord, error) { return nil, nil }
func (m *mockDataStore) GetOCRByJobID(jobID string) (*types.JobRecord, error)       { return nil, nil }

type mockFaceMatch struct{}

func (mfm *mockFaceMatch) CalcFaceMatchScore(payload types.FaceMatchPayload) (int, error) {
	return 45, nil
}

type mockOCR struct{}

func (mfm *mockOCR) PerformOCR(payload types.OCRPayload) (*types.OCRResponse, error) {
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

type mockUuid struct{}

func (u *mockUuid) New() string {
	return "new-uuid"
}

type mockTaskQueue struct{}

func (tq *mockTaskQueue) PushJobOnQueue(payload []byte) error             { return nil }
func (tq *mockTaskQueue) PullJobFromQueue() (<-chan amqp.Delivery, error) { return nil, nil }

func TestValidateImage(t *testing.T) {
	tt := []struct {
		name     string
		payload  types.FaceMatchPayload
		clientID int
		expErr   error
	}{
		{
			name: "not a face image for first image id",
			payload: types.FaceMatchPayload{
				Image1: "abc",
				Image2: "xyz",
			},
			clientID: 1,
			expErr:   ErrNotFaceImg,
		},
		{
			name: "not a face image for second image id",
			payload: types.FaceMatchPayload{
				Image1: "def",
				Image2: "pqr",
			},
			clientID: 2,
			expErr:   ErrNotFaceImg,
		},
		{
			name: "non-existent image for first image id",
			payload: types.FaceMatchPayload{
				Image1: "ert",
				Image2: "pqr",
			},
			clientID: 1,
			expErr:   ErrInvalidImgId,
		},
		{
			name: "non-existent image for second image id",
			payload: types.FaceMatchPayload{
				Image1: "def",
				Image2: "jnk",
			},
			clientID: 1,
			expErr:   ErrInvalidImgId,
		},
		{
			name: "different client id for both images",
			payload: types.FaceMatchPayload{
				Image1: "xyz",
				Image2: "def",
			},
			clientID: 1,
			expErr:   ErrInvalidImgId,
		},
		{
			name: "different client id for client and image",
			payload: types.FaceMatchPayload{
				Image1: "asdf",
				Image2: "lkjh",
			},
			clientID: 2,
			expErr:   ErrInvalidImgId,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			service := &Service{
				dataStore: &mockDataStore{},
			}

			err := service.ValidateImage(tc.payload, tc.clientID)
			if err != tc.expErr {
				t.Errorf("Expected %q but got %q", tc.expErr, err)
			}
		})
	}
}

func TestCalcFaceMatchScore(t *testing.T) {
	tt := []struct {
		name    string
		payload types.FaceMatchPayload
		expOut  int
		expErr  error
	}{
		{
			name: "only case",
			payload: types.FaceMatchPayload{
				Image1: "abc",
				Image2: "xyz",
			},
			expOut: 45,
			expErr: nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			service := &Service{
				dataStore: &mockDataStore{},
				faceMatch: &mockFaceMatch{},
			}
			clientID := 1
			score, err := service.CalcAndSaveFaceMatchScore(tc.payload, clientID)
			if tc.expErr != nil {
				if err == nil {
					t.Fatalf("Expected error but got nil")
				}

				if !errors.Is(tc.expErr, err) {
					t.Errorf("Expected error %q but got %q", tc.expErr, err)
				}
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if score != tc.expOut {
				t.Errorf("Expected %q but got %q", tc.expOut, score)
			}
		})
	}
}

func TestValidateImageOCR(t *testing.T) {
	tt := []struct {
		name     string
		payload  types.OCRPayload
		clientID int
		expErr   error
	}{
		{
			name: "not an id_card image for image id",
			payload: types.OCRPayload{
				Image: "cvbas",
			},
			clientID: 2,
			expErr:   ErrNotIDCardImg,
		},
		{
			name: "non-existent image for first image id",
			payload: types.OCRPayload{
				Image: "ert",
			},
			clientID: 2,
			expErr:   ErrInvalidImgId,
		},
		{
			name: "different client id for image and client",
			payload: types.OCRPayload{
				Image: "cvbasrt",
			},
			clientID: 4,
			expErr:   ErrInvalidImgId,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			service := &Service{
				dataStore: &mockDataStore{},
			}

			err := service.ValidateImageOCR(tc.payload, tc.clientID)
			if err != tc.expErr {
				t.Errorf("Expected %q but got %q", tc.expErr, err)
			}
		})
	}
}

func TestPerformOCR(t *testing.T) {
	tt := []struct {
		name    string
		payload types.OCRPayload
		expOut  *types.OCRResponse
		expErr  error
	}{
		{
			name: "only case",
			payload: types.OCRPayload{
				Image: "ac",
			},
			expOut: &types.OCRResponse{
				Name:      "John Adams",
				Gender:    "Male",
				DOB:       "1990-01-24",
				IdNumber:  "1234-1234-1234",
				AddrLine1: "A2, 201, Amar Villa",
				AddrLine2: "MG Road, Pune",
				Pincode:   "411004",
			},
			expErr: nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			service := &Service{
				dataStore:  &mockDataStore{},
				ocrService: &mockOCR{},
			}

			clientID := 1
			resp, err := service.PerformAndSaveOCR(tc.payload, clientID)
			if tc.expErr != nil {
				if err == nil {
					t.Fatalf("Expected error but got nil")
				}

				if !errors.Is(tc.expErr, err) {
					t.Errorf("Expected error %q but got %q", tc.expErr, err)
				}
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !reflect.DeepEqual(resp, tc.expOut) {
				t.Errorf("Expected %q but got %q", tc.expOut, resp)
			}
		})
	}
}

func TestPerformFaceMatchAsync(t *testing.T) {
	tt := []struct {
		name     string
		payload  types.FaceMatchPayload
		clientID int
		expOut   string
		expErr   error
	}{
		{
			name: "not a face image for first image id",
			payload: types.FaceMatchPayload{
				Image1: "abc",
				Image2: "xyz",
			},
			clientID: 1,
			expErr:   ErrNotFaceImg,
		},
		{
			name: "not a face image for second image id",
			payload: types.FaceMatchPayload{
				Image1: "def",
				Image2: "pqr",
			},
			clientID: 2,
			expErr:   ErrNotFaceImg,
		},
		{
			name: "non-existent image for first image id",
			payload: types.FaceMatchPayload{
				Image1: "ert",
				Image2: "pqr",
			},
			clientID: 1,
			expErr:   ErrInvalidImgId,
		},
		{
			name: "non-existent image for second image id",
			payload: types.FaceMatchPayload{
				Image1: "def",
				Image2: "jnk",
			},
			clientID: 1,
			expErr:   ErrInvalidImgId,
		},
		{
			name: "different client id for both images",
			payload: types.FaceMatchPayload{
				Image1: "xyz",
				Image2: "def",
			},
			clientID: 1,
			expErr:   ErrInvalidImgId,
		},
		{
			name: "different client id for client and image",
			payload: types.FaceMatchPayload{
				Image1: "asdf",
				Image2: "lkjh",
			},
			clientID: 2,
			expErr:   ErrInvalidImgId,
		},
		{
			name: "valid case",
			payload: types.FaceMatchPayload{
				Image1: "asdf",
				Image2: "lkjh",
			},
			clientID: 3,
			expOut:   "new-uuid",
			expErr:   nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			service := &Service{
				dataStore: &mockDataStore{},
				faceMatch: &mockFaceMatch{},
				uuid:      &mockUuid{},
				queue:     &mockTaskQueue{},
			}

			id, err := service.PerformFaceMatchAsync(tc.payload, tc.clientID)
			if tc.expErr != nil {
				if err == nil {
					t.Fatalf("Expected error but got nil")
				}

				if !errors.Is(tc.expErr, err) {
					t.Errorf("Expected error %q but got %q", tc.expErr, err)
				}

				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if id != tc.expOut {
				t.Errorf("Expected %q but got %q", tc.expOut, id)
			}
		})
	}
}

func TestPerformOCRAsync(t *testing.T) {
	tt := []struct {
		name     string
		payload  types.OCRPayload
		clientID int
		expOut   string
		expErr   error
	}{
		{
			name: "not an id_card image for image id",
			payload: types.OCRPayload{
				Image: "cvbas",
			},
			clientID: 2,
			expErr:   ErrNotIDCardImg,
		},
		{
			name: "non-existent image for first image id",
			payload: types.OCRPayload{
				Image: "ert",
			},
			clientID: 2,
			expErr:   ErrInvalidImgId,
		},
		{
			name: "different client id for image and client",
			payload: types.OCRPayload{
				Image: "cvbasrt",
			},
			clientID: 4,
			expErr:   ErrInvalidImgId,
		},
		{
			name: "only case",
			payload: types.OCRPayload{
				Image: "ac",
			},
			clientID: 3,
			expOut:   "new-uuid",
			expErr:   nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			service := &Service{
				dataStore:  &mockDataStore{},
				ocrService: &mockOCR{},
				uuid:       &mockUuid{},
				queue:      &mockTaskQueue{},
			}

			id, err := service.PerformOCRAsync(tc.payload, tc.clientID)
			if tc.expErr != nil {
				if err == nil {
					t.Fatalf("Expected error but got nil")
				}

				if !errors.Is(tc.expErr, err) {
					t.Errorf("Expected error %q but got %q", tc.expErr, err)
				}

				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if id != tc.expOut {
				t.Errorf("Expected %q but got %q", tc.expOut, id)
			}
		})
	}
}
