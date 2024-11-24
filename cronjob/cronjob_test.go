package cronjob

import (
	"testing"

	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
)

type mockCronJobStore struct {
	counter int
}

func (mst *mockCronJobStore) GetReportData(date string) ([]*types.ClientReport, error) {
	mst.counter++
	return []*types.ClientReport{
		{
			ClientID:          "1",
			Name:              "client1",
			Plan:              "basic",
			Date:              "2024-11-22",
			TotalFaceMatch:    "5",
			TotalOcr:          "2",
			TotalImgStorageMB: "10",
			TotalAPIUsageCost: "12",
			TotalStorageCost:  "10",
		},
		{
			ClientID:          "2",
			Name:              "client2",
			Plan:              "basic",
			Date:              "2024-11-22",
			TotalFaceMatch:    "7",
			TotalOcr:          "4",
			TotalImgStorageMB: "12",
			TotalAPIUsageCost: "14",
			TotalStorageCost:  "12",
		},
	}, nil
}

type mockCronJobService struct {
	counter int
}

func (mse *mockCronJobService) PrepareCSV([]*types.ClientReport) ([]byte, error) {
	mse.counter++
	return nil, nil
}

type mockCronJobFileStore struct {
	counter int
}

func (mfs *mockCronJobFileStore) SaveFile(file *types.FileUpload) error {
	mfs.counter++
	return nil
}

func (mfs *mockCronJobFileStore) GetFile(filePath string) ([]byte, error) {
	return nil, nil
}

func TestCalcDailyReport(t *testing.T) {
	// call the method
	mockService := &mockCronJobService{}
	mockDataStore := &mockCronJobStore{}
	mockFileStore := &mockCronJobFileStore{}
	cj := NewCronJob(mockDataStore, mockFileStore, mockService)
	cj.CalcDailyReport()

	// check if the required call stack was followed
	if mockService.counter != 1 {
		t.Errorf("Expected mock service counter to be %d but got %d", 1, mockService.counter)
	}
	if mockDataStore.counter != 1 {
		t.Errorf("Expected mock data store counter to be %d but got %d", 1, mockDataStore.counter)
	}
	if mockFileStore.counter != 1 {
		t.Errorf("Expected mock file store counter to be %d but got %d", 1, mockFileStore.counter)
	}
}
