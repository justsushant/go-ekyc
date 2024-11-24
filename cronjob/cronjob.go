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
		Name:    fmt.Sprintf("reports/%s", strings.ReplaceAll(currentDate, "-", "")),
		Content: bytes.NewReader(csvBytes),
		Size:    int64(len(csvBytes)),
		Headers: map[string]string{
			"Content-Type": "text/csv",
		},
	}
	c.fileStore.SaveFile(file)

}
