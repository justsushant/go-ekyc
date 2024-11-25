package cronjob

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/store"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
)

type CronJob struct {
	service   CronJobServiceInterface
	dataStore store.CronJobDataStore
	fileStore store.FileStore
}

func NewCronJob(dataStore store.CronJobDataStore, fileStore store.FileStore, service CronJobServiceInterface) *CronJob {
	return &CronJob{
		dataStore: dataStore,
		fileStore: fileStore,
		service:   service,
	}
}

func (c *CronJob) CalcDailyReport() {
	// calc the required data from store layer
	currentDate := time.Now().Format("2006-01-02")
	data, err := c.dataStore.GetReportData(currentDate)
	if err != nil {
		log.Printf("Error while fetching data from store for %s report: %s\n", currentDate, err.Error())
	}

	// convert into csv file
	csvBytes, err := c.service.PrepareCSV(data)
	if err != nil {
		log.Printf("Error while preparing csv file for date %s: %s\n", currentDate, err.Error())
	}

	// save to file store
	file := &types.FileUpload{
		Name:    c.getDailyReportPath(strings.ReplaceAll(currentDate, "-", "")),
		Content: bytes.NewReader(csvBytes),
		Size:    int64(len(csvBytes)),
		Headers: map[string]string{
			"Content-Type": "text/csv",
		},
	}
	c.fileStore.SaveFile(file)
}

func (c *CronJob) CalcMonthlyReport(currentTime time.Time) {
	// calc the required data from store layer
	currentMonth := currentTime.Month()
	currentYear := currentTime.Year()
	data, err := c.dataStore.GetMonthlyReport(int(currentMonth), currentYear)
	if err != nil {
		log.Printf("Error while fetching data from store for month %d report: %s\n", currentMonth, err.Error())
	}

	for _, d := range data {
		// extract client id
		var clientID string
		if len(d) != 0 {
			clientID = d[0].ClientID
		}

		// convert into csv file
		csvBytes, err := c.service.PrepareCSV(d)
		if err != nil {
			log.Printf("Error while preparing csv file for month %d: %s\n", currentMonth, err.Error())
		}

		// save to file store
		file := &types.FileUpload{
			Name:    c.getMonthlyReportPath(clientID, int(currentMonth), currentYear),
			Content: bytes.NewReader(csvBytes),
			Size:    int64(len(csvBytes)),
			Headers: map[string]string{
				"Content-Type": "text/csv",
			},
		}
		c.fileStore.SaveFile(file)
	}
}

func (c *CronJob) getDailyReportPath(date string) string {
	return fmt.Sprintf("reports/daily/%s", strings.ReplaceAll(date, "-", ""))
}

func (c *CronJob) getMonthlyReportPath(clientID string, month, year int) string {
	return fmt.Sprintf("reports/monthly/%s-%d%d", clientID, month, year)
}
