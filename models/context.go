package models

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB(connectionString string) error {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}

	DB = db
	return err
}
