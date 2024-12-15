package main

import (
	"flag"
	"log"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/config"
)

const MIGRATION_FILES_PATH = "file://db/migration"

var (
	mode    string
	version int
	isForce bool
)

// TODO: Implement the force option
// TODO: Implement the up and down with steps, by default only one step up or down

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatalf("Error while config init: %v", err)
	}

	// set the flags
	flag.StringVar(&mode, "m", "", "migrate [up|down|status]")
	flag.IntVar(&version, "s", 0, "migration steps to apply")
	flag.IntVar(&version, "v", 0, "which migration version to apply")
	flag.BoolVar(&isForce, "f", false, "force migration or not")
	flag.Parse()

	// create migration using migration file path and db conn string
	mig, err := migrate.New(MIGRATION_FILES_PATH, cfg.DbDsn)
	if err != nil {
		log.Fatalf("Error occured while creating migration: %v\n", err)
	}

	// acting on the supplied migration mode
	switch strings.ToLower(mode) {
	case "up":
		if version == 0 {
			err = mig.Steps(1)
			if err != nil && err.Error() != "no change" {
				log.Fatalf("Error occured while applying migrations: %v\n", err)
			}
			log.Printf("Migration applied successfully\n")
			version, dirty, err := mig.Version()
			if err != nil {
				log.Fatalf("Error occured while getting migration status: %v\n", err)
			}
			log.Printf("Current version: %d, Dirty: %v\n", version, dirty)
			return
		} else {
			for range version {
				err = mig.Steps(1)
				if err != nil && err.Error() != "no change" {
					log.Fatalf("Error occured while applying migrations: %v\n", err)
				}
			}
			log.Printf("Migration applied successfully\n")
			version, dirty, err := mig.Version()
			if err != nil {
				log.Fatalf("Error occured while getting migration status: %v\n", err)
			}
			log.Printf("Current version: %d, Dirty: %v\n", version, dirty)
			return
		}
	case "down":
		if version == 0 {
			err = mig.Steps(-1)
			if err != nil && err.Error() != "no change" {
				log.Fatalf("Error occured while applying migrations: %v\n", err)
			}
			log.Printf("Rollback applied successfully\n")
			version, dirty, err := mig.Version()
			if err != nil {
				log.Fatalf("Error occured while getting migration status: %v\n", err)
			}
			log.Printf("Current version: %d, Dirty: %v\n", version, dirty)
			return
		} else {
			for range version {
				err = mig.Steps(-1)
				if err != nil && err.Error() != "no change" {
					log.Fatalf("Error occured while applying migrations: %v\n", err)
				}
			}
			log.Printf("Rollback applied successfully\n")
			version, dirty, err := mig.Version()
			if err != nil {
				log.Fatalf("Error occured while getting migration status: %v\n", err)
			}
			log.Printf("Current version: %d, Dirty: %v\n", version, dirty)
			return
		}
	case "status":
		version, dirty, err := mig.Version()
		if err != nil {
			log.Fatalf("Error occured while getting migration status: %v\n", err)
		}
		log.Printf("Current version: %d, Dirty: %v\n", version, dirty)
		return
	}

	// if specific migration version was passed and force flag was not passed
	if version != 0 && !isForce {
		err := mig.Migrate(uint(version))
		if err != nil {
			log.Fatalf("Error occured while applying migration: %v", err)
		}
		log.Printf("Migration version %d applied successfully\n", version)
	}

	// if specific migration version & force flag was passed
	if version != 0 && isForce {
		err := mig.Force(version)
		if err != nil {
			log.Fatalf("Error occured while applying migration: %v", err)
		}
		log.Printf("Migration version %d applied successfully\n", version)
	}

}
