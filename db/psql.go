package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func NewPsqlClient(dsn string) *sql.DB {
	// connect to database
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// ping to database
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("PostgreSQL client connected")
	return db
}
