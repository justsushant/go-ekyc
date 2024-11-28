package cronjob

import (
	"bytes"
	"errors"
	"log"

	"github.com/gocarina/gocsv"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
)

type CronJobServiceManager interface {
	PrepareCSV(data interface{}) ([]byte, error)
}

var ErrMissingReports = errors.New("no reports found")

type CronJobService struct{}

func NewCronJobService() *CronJobService {
	return &CronJobService{}
}

func (s *CronJobService) PrepareCSV(data interface{}) ([]byte, error) {
	var buf bytes.Buffer

	switch d := data.(type) {
	case []*types.ClientReport:
		// length check
		if len(d) == 0 || data == nil {
			return nil, ErrMissingReports
		}

		// writing the data to buffer
		if err := gocsv.Marshal(&d, &buf); err != nil {
			log.Printf("Error while writing CSV to buffer: %s\n", err.Error())
			return nil, err
		}
	case []*types.ClientReportMonthly:
		// length check
		if len(d) == 0 || data == nil {
			return nil, ErrMissingReports
		}

		// writing the data to buffer
		if err := gocsv.Marshal(&d, &buf); err != nil {
			log.Printf("Error while writing CSV to buffer: %s\n", err.Error())
			return nil, err
		}
	}

	return buf.Bytes(), nil
}
