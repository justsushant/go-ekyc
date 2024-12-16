package cronjob

import (
	"testing"
	"time"

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
func (mst *mockCronJobStore) GetMonthlyReport(currentMonth, currentYear int) ([][]*types.ClientReportMonthly, error) {
	mst.counter++
	return [][]*types.ClientReportMonthly{
		{
			{
				ClientID:          "1",
				Date:              "2024-11-22",
				TotalFaceMatch:    "5",
				TotalOcr:          "2",
				TotalImgStorageMB: "10",
				TotalAPIUsageCost: "12",
				TotalStorageCost:  "10",
			},
			{
				ClientID:          "1",
				Date:              "2024-11-23",
				TotalFaceMatch:    "7",
				TotalOcr:          "4",
				TotalImgStorageMB: "12",
				TotalAPIUsageCost: "14",
				TotalStorageCost:  "12",
			},
		},
		{
			{
				ClientID:          "2",
				Date:              "2024-11-22",
				TotalFaceMatch:    "7",
				TotalOcr:          "4",
				TotalImgStorageMB: "12",
				TotalAPIUsageCost: "14",
				TotalStorageCost:  "12",
			},
			{
				ClientID:          "2",
				Date:              "2024-11-23",
				TotalFaceMatch:    "5",
				TotalOcr:          "2",
				TotalImgStorageMB: "10",
				TotalAPIUsageCost: "12",
				TotalStorageCost:  "10",
			},
		},
	}, nil
}

type mockCronJobService struct {
	counter int
}

func (mse *mockCronJobService) PrepareCSV(data interface{}) ([]byte, error) {
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
	cj := &CronJob{
		service:   mockService,
		db:        mockDataStore,
		fileStore: mockFileStore,
	}
	cj.CalcDailyReport(time.Date(2024, time.November, 1, 0, 0, 0, 0, time.UTC))

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

func TestCalcMonthlyReport(t *testing.T) {
	tt := []struct {
		name                string
		time                time.Time
		expDStoreCallCount  int
		expFStoreCallCount  int
		expServiceCallCount int
	}{
		{
			name:                "monthly report with two clients",
			time:                time.Date(2024, time.November, 1, 0, 0, 0, 0, time.UTC),
			expDStoreCallCount:  1,
			expFStoreCallCount:  2,
			expServiceCallCount: 2,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// call the method
			mockService := &mockCronJobService{}
			mockDataStore := &mockCronJobStore{}
			mockFileStore := &mockCronJobFileStore{}
			cj := &CronJob{
				service:   mockService,
				db:        mockDataStore,
				fileStore: mockFileStore,
			}
			cj.CalcMonthlyReport(tc.time)

			// check if the required call stack was followed
			if mockService.counter != tc.expServiceCallCount {
				t.Errorf("Expected mock service counter to be %d but got %d", tc.expServiceCallCount, mockService.counter)
			}
			if mockDataStore.counter != tc.expDStoreCallCount {
				t.Errorf("Expected mock data store counter to be %d but got %d", tc.expDStoreCallCount, mockDataStore.counter)
			}
			if mockFileStore.counter != tc.expFStoreCallCount {
				t.Errorf("Expected mock file store counter to be %d but got %d", tc.expFStoreCallCount, mockFileStore.counter)
			}
		})
	}
}
