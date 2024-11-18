package service

import (
	"errors"
	"testing"

	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
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

	return nil, nil
}

type mockFaceMatch struct{}

func (mfm *mockFaceMatch) CalcFaceMatchScore(payload types.FaceMatchPayload) (int, error) {
	return 45, nil
}

func TestValidateImage(t *testing.T) {
	tt := []struct {
		name    string
		payload types.FaceMatchPayload
		expErr  error
	}{
		{
			name: "not a face image for first image id",
			payload: types.FaceMatchPayload{
				ImageID1: "abc",
				ImageID2: "xyz",
			},
			expErr: ErrNotFaceImg,
		},
		{
			name: "not a face image for second image id",
			payload: types.FaceMatchPayload{
				ImageID1: "def",
				ImageID2: "pqr",
			},
			expErr: ErrNotFaceImg,
		},
		{
			name: "non-existent image for first image id",
			payload: types.FaceMatchPayload{
				ImageID1: "ert",
				ImageID2: "pqr",
			},
			expErr: ErrInvalidImgId,
		},
		{
			name: "non-existent image for second image id",
			payload: types.FaceMatchPayload{
				ImageID1: "def",
				ImageID2: "jnk",
			},
			expErr: ErrInvalidImgId,
		},
		{
			name: "different client id for both images",
			payload: types.FaceMatchPayload{
				ImageID1: "xyz",
				ImageID2: "def",
			},
			expErr: ErrInvalidImgId,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			service := &Service{
				dataStore: &mockDataStore{},
			}

			err := service.ValidateImage(tc.payload)
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
				ImageID1: "abc",
				ImageID2: "xyz",
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

			score, err := service.CalcFaceMatchScore(tc.payload)
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
