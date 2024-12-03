package db

import (
	"database/sql"
	"log"
	"net/url"

	_ "github.com/lib/pq"
)

type PostgresConn struct {
	Endpoint string
	User     string
	Password string
	Ssl      string
	Db       string
}

func NewPsqlClient(conn *PostgresConn) *sql.DB {
	// making psql conn string
	dsn := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(conn.User, conn.Password),
		Host:     conn.Endpoint,
		Path:     conn.Db,
		RawQuery: "sslmode=disable",
	}

	// connect to database
	db, err := sql.Open("postgres", dsn.String())
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
