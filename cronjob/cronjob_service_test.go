package cronjob

import (
	"errors"
	"testing"

	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
)

func TestPrepareCSV(t *testing.T) {
	tt := []struct {
		name    string
		data    []*types.ClientReport
		expResp string
		expErr  error
	}{
		{
			name:    "zero report objects",
			data:    []*types.ClientReport{},
			expResp: "",
			expErr:  ErrMissingReports,
		},
		{
			name: "one report object",
			data: []*types.ClientReport{
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
			},
			expResp: "client_id,name,plan,date,total_face_match_for_day,total_ocr_for_day,total_image_storage_in_mb,api_usage_cost_usd,storage_cost_usd\n1,client1,basic,2024-11-22,5,2,10,12,10\n",
			expErr:  nil,
		},
		{
			name: "two report objects",
			data: []*types.ClientReport{
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
			},
			expResp: "client_id,name,plan,date,total_face_match_for_day,total_ocr_for_day,total_image_storage_in_mb,api_usage_cost_usd,storage_cost_usd\n1,client1,basic,2024-11-22,5,2,10,12,10\n2,client2,basic,2024-11-22,7,4,12,14,12\n",
			expErr:  nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			service := NewCronJobService()
			csvBytes, err := service.PrepareCSV(tc.data)

			if tc.expErr != nil {
				if err == nil {
					t.Errorf("Expected error but didn't got one\n")
				}

				if !errors.Is(err, tc.expErr) {
					t.Errorf("Expected error %s but got %s\n", tc.expErr, err)
				}

				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %s", err.Error())
			}

			if string(csvBytes) != tc.expResp {
				t.Errorf("Expected %q but got %q", tc.expResp, string(csvBytes))
			}
		})
	}
}

func TestPrepareCSVForMonthlyReport(t *testing.T) {
	tt := []struct {
		name    string
		data    []*types.ClientReportMonthly
		expResp string
		expErr  error
	}{
		{
			name: "two report objects",
			data: []*types.ClientReportMonthly{
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
			expResp: "client_id,date,total_face_match_for_day,total_ocr_for_day,total_image_storage_in_mb,api_usage_cost_usd,storage_cost_usd\n1,2024-11-22,5,2,10,12,10\n1,2024-11-23,7,4,12,14,12\n",
			expErr:  nil,
		},
		{
			name: "one report object",
			data: []*types.ClientReportMonthly{
				{
					ClientID:          "1",
					Date:              "2024-11-22",
					TotalFaceMatch:    "5",
					TotalOcr:          "2",
					TotalImgStorageMB: "10",
					TotalAPIUsageCost: "12",
					TotalStorageCost:  "10",
				},
			},
			expResp: "client_id,date,total_face_match_for_day,total_ocr_for_day,total_image_storage_in_mb,api_usage_cost_usd,storage_cost_usd\n1,2024-11-22,5,2,10,12,10\n",
			expErr:  nil,
		},
		{
			name:    "zero report objects",
			data:    []*types.ClientReportMonthly{},
			expResp: "",
			expErr:  ErrMissingReports,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			service := NewCronJobService()
			csvBytes, err := service.PrepareCSV(tc.data)

			if tc.expErr != nil {
				if err == nil {
					t.Errorf("Expected error but didn't got one\n")
				}

				if !errors.Is(err, tc.expErr) {
					t.Errorf("Expected error %s but got %s\n", tc.expErr, err)
				}

				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %s", err.Error())
			}

			if string(csvBytes) != tc.expResp {
				t.Errorf("Expected %q but got %q", tc.expResp, string(csvBytes))
			}
		})
	}
}
