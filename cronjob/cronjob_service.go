package cronjob

import "github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"

// TODO: change the name of interface
type CronJobServiceInterface interface {
	PrepareCSV([]*types.ClientReport) ([]byte, error)
}

type CronJobService struct{}

func NewCronJobService() *CronJobService {
	return &CronJobService{}
}

func (s *CronJobService) PrepareCSV(*[]types.ClientReport) ([]byte, error) {
	return nil, nil
}
