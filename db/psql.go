package db

import (
	"database/sql"
	"fmt"
	"log"

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
	// dsn := url.URL{
	// 	Scheme:   "postgres",
	// 	User:     url.UserPassword(conn.User, conn.Password),
	// 	Host:     conn.Endpoint,
	// 	Path:     conn.Db,
	// 	RawQuery: "sslmode=disable",
	// }

	dsn := fmt.Sprintf("postgres://%s:%s@database:5432/ekyc_db?sslmode=disable", conn.User, conn.Password)

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
