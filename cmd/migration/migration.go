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
	migVer  int
	isForce bool
)

// TODO: Implement the force option
// TODO: Implement the up and down with steps, by default only one step up or down

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalf("Error while config init: %w", err)
	}

	// set the flags
	flag.StringVar(&mode, "m", "", "migrate [up|down|status]")
	flag.IntVar(&migVer, "v", 0, "which migration version to apply")
	flag.BoolVar(&isForce, "f", false, "force migration or not")
	flag.Parse()

	// extract database conn string
	pgDsn := cfg.DbDsn
	if pgDsn == "" {
		log.Fatalf("Error: PostgreSQL dsn not found\n")
	}

	// create migration
	mig, err := migrate.New(MIGRATION_FILES_PATH, pgDsn)
	if err != nil {
		log.Fatalf("Error occured while creating migration: %v\n", err)
	}

	// acting on the supplied migration mode
	if mode != "" {
		switch strings.ToLower(mode) {
		case "up":
			err = mig.Up()
			if err != nil && err.Error() != "no change" {
				log.Fatalf("Error occured while applying migrations: %v\n", err)
			}
			log.Printf("Migrations applied successfully\n")
		case "down":
			err = mig.Down()
			if err != nil && err.Error() != "no change" {
				log.Fatalf("Error occured while rolling back migrations: %v\n", err)
			}
			log.Printf("Rollback applied successfully\n")
		case "status":
			version, dirty, err := mig.Version()
			if err != nil {
				log.Fatalf("Error occured while getting migration status: %v\n", err)
			}
			log.Printf("Current version: %d, Dirty: %v\n", version, dirty)
		default:
			log.Fatal("Unknown command")
		}
	}

	// if specific migration version was passed
	if migVer != 0 {
		err := mig.Migrate(uint(migVer))
		if err != nil {
			log.Fatalf("Error occured while applying migration: %v", err)
		}
		log.Printf("Migration version %d applied successfully\n", migVer)
	}
}
