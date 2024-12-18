package cronjob

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/store"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
	"github.com/robfig/cron/v3"
)

type CronJob struct {
	service   CronJobServiceManager
	db        store.CronJobDataStore
	fileStore store.FileStore
	Cron      *cron.Cron
}

type CronJobConfig struct {
	ServiceManager CronJobServiceManager
	DataStore      store.CronJobDataStore
	FileStore      store.FileStore
	Cron           *cron.Cron
}

// instantiate a CronJob struct
func New(cronjobConfig *CronJobConfig) *CronJob {
	return &CronJob{
		service:   cronjobConfig.ServiceManager,
		db:        cronjobConfig.DataStore,
		fileStore: cronjobConfig.FileStore,
		Cron:      cronjobConfig.Cron,
	}
}

func (c *CronJob) CalcDailyReport(currentTime time.Time) {
	// calc the required data from store layer
	currentDateString := currentTime.Format("2006-01-02")
	data, err := c.db.GetReportData(currentDateString)
	if err != nil {
		log.Printf("Error while fetching data from store for %s report: %s\n", currentDateString, err.Error())
		return
	}

	// convert into csv file
	csvBytes, err := c.service.PrepareCSV(data)
	if err != nil {
		log.Printf("Error while preparing csv file for date %s: %s\n", currentDateString, err.Error())
		return
	}

	// save to file store
	file := &types.FileUpload{
		Name:    c.getDailyReportPath(strings.ReplaceAll(currentDateString, "-", "")),
		Content: bytes.NewReader(csvBytes),
		Size:    int64(len(csvBytes)),
		Headers: map[string]string{
			"Content-Type": "text/csv",
		},
	}
	err = c.fileStore.SaveFile(file)
	if err != nil {
		log.Printf("Error while saving csv file for date %s: %s\n", currentDateString, err.Error())
		return
	}

	log.Printf("Daily report cronjob executed successfully at %s", time.Now().String())
}

func (c *CronJob) CalcMonthlyReport(currentTime time.Time) {
	// calc the required data from store layer
	currentMonth := currentTime.Month()
	currentYear := currentTime.Year()
	data, err := c.db.GetMonthlyReport(int(currentMonth), currentYear)
	if err != nil {
		log.Printf("Error while fetching data from store for month %d report: %s\n", currentMonth, err.Error())
		return
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
			log.Printf("Error while preparing csv file of clientID %s for month %d: %s\n", clientID, currentMonth, err.Error())
			continue
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
		err = c.fileStore.SaveFile(file)
		if err != nil {
			log.Printf("Error while saving csv file for month-year %d-%d: %s\n", int(currentMonth), currentYear, err.Error())
			return
		}
	}

	log.Printf("Monthly report cronjob executed successfully at %s", time.Now().String())
}

func (c *CronJob) getDailyReportPath(date string) string {
	return fmt.Sprintf("reports/daily/%s", strings.ReplaceAll(date, "-", ""))
}

func (c *CronJob) getMonthlyReportPath(clientID string, month, year int) string {
	return fmt.Sprintf("reports/monthly/%s-%d%d", clientID, month, year)
}
