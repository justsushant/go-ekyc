package store

import "github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"

type CronJobDataStore interface {
	GetReportData(date string) ([]*types.ClientReport, error)
	GetMonthlyReport(currentMonth, currentYear int) ([]*types.ClientReportMonthly, error)
}
