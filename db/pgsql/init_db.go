package pgsql

import (
	"database/sql"
	"fmt"
	"urlshortener/app/starter"

	_ "github.com/lib/pq"
)

func InitDB(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return db, fmt.Errorf("error connection db: %v", err)
	}

	return db, nil
}

func InitStorers(db *sql.DB) *starter.Storers {
	u := NewUserStore(db)
	l := NewLinkStorer(db)
	lt := NewLinkTransitStorer(db)

	return starter.NewStorers(u, l, lt)
}
