package cronjob

import (
	"bytes"
	"errors"
	"log"

	"github.com/gocarina/gocsv"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
)

// TODO: change the name of interface
type CronJobServiceInterface interface {
	PrepareCSV([]*types.ClientReport) ([]byte, error)
}

var ErrMissingReports = errors.New("no reports found")

type CronJobService struct{}

func NewCronJobService() *CronJobService {
	return &CronJobService{}
}

func (s *CronJobService) PrepareCSV(data []*types.ClientReport) ([]byte, error) {
	// length check
	if len(data) == 0 {
		return nil, ErrMissingReports
	}

	// buffer to hold the csv
	var buf bytes.Buffer

	// writing the data to buffer
	if err := gocsv.Marshal(&data, &buf); err != nil {
		log.Printf("Error while writing CSV to buffer: %s\n", err.Error())
		return nil, err
	}

	return buf.Bytes(), nil
}
