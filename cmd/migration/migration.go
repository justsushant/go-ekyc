package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/config"
)

const MIGRATION_FILES_PATH = "file://migration"

var (
	mode    string
	migVer  int
	isForce bool
)

// TODO: Implement the migration version and force option

func main() {
	// set the flags
	flag.StringVar(&mode, "m", "", "migrate [up|down|status]")
	flag.IntVar(&migVer, "v", 0, "which migration version to apply")
	flag.BoolVar(&isForce, "f", false, "force migration or not")
	flag.Parse()

	// extract database conn string
	pgDsn := config.Envs.DB_Dsn
	if pgDsn == "" {
		panic("postgresql dsn not found")
	}

	// create migration
	mig, err := migrate.New(MIGRATION_FILES_PATH, pgDsn)
	if err != nil {
		panic(fmt.Sprintf("error occured while creating migration: %v", err))
	}

	// acting on the supplied migration mode
	switch strings.ToLower(mode) {
	case "up":
		err = mig.Up()
		if err != nil && err.Error() != "no change" {
			log.Fatalf("Error applying migrations: %v\n", err)
		}
		fmt.Println("Migrations applied successfully")
	case "down":
		err = mig.Down()
		if err != nil && err.Error() != "no change" {
			log.Fatalf("Error rolling back migrations: %v\n", err)
		}
		fmt.Println("Rollback applied successfully")
	case "status":
		version, dirty, err := mig.Version()
		if err != nil {
			log.Fatalf("Error getting migration status: %v\n", err)
		}
		fmt.Printf("Current version: %d, Dirty: %v\n", version, dirty)
	default:
		log.Fatal("Unknown command")
	}
}
