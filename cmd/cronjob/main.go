package main

import (
	"log"
	"time"

	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/config"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/cronjob"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/db"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/service"
	"github.com/robfig/cron/v3"
)

func main() {
	// load configs
	cfg, err := config.Init()
	if err != nil {
		log.Fatalf("Error while config init: %v", err)
	}

	// get psql store
	psqlStore := cronjob.NewPsqlCrobJobStore(cfg.DbDsn)

	// get minio store
	minioConn := &db.MinioConn{
		Endpoint: cfg.MinioEndpoint,
		User:     cfg.MinioUser,
		Password: cfg.MinioPassword,
		Ssl:      cfg.MinioSSL,
	}
	minioStore := service.NewMinioStore(minioConn, cfg.MinioBucket)

	// get cronjob service
	service := cronjob.NewCronJobService()

	// start the cronjob
	c := cronjob.NewCronJob(psqlStore, minioStore, service, cron.New())

	// add the schedules
	currentTime := time.Now()
	_, err = c.Cron.AddFunc("0 1 * * *", func() { // every day at 1AM
		c.CalcDailyReport(currentTime)
	})
	if err != nil {
		log.Println("Error scheduling the job:", err.Error())
		return
	}

	_, err = c.Cron.AddFunc("0 1 1 * *", func() { // first day of every month at 1AM
		c.CalcMonthlyReport(currentTime)
	})
	if err != nil {
		log.Println("Error scheduling monthly job:", err.Error())
		return
	}

	// start the job
	c.Cron.Start()

	// for testing purposes
	c.CalcDailyReport(currentTime)
	c.CalcMonthlyReport(currentTime)

	select {} // to hold the program from exiting
}
