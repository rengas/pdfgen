package dbutils

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/rengas/pdfgen/pkg/retry"
	"log"
	"time"
)

func OpenPostgres(connString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	if err = retry.Retry(5, time.Second, 3, db.Ping); err != nil {
		return nil, fmt.Errorf("failed to ping postgres server")
	}

	return db, nil
}
func MustOpenPostgres(connString string) *sql.DB {
	db, err := OpenPostgres(connString)
	if err != nil {
		log.Fatal("failed to connect to posgres")
	}
	return db
}
