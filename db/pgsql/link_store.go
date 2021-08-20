package pgsql

import (
	"database/sql"
	"fmt"
	"urlshortener/app/services/link"
)

var _ link.LinkStorer = &LinkStore{}

type LinkStore struct {
	db *sql.DB
}

func NewLinkStorer(db *sql.DB) *LinkStore {
	return &LinkStore{db: db}
}

func (l *LinkStore) Insert(link link.Link) error {
	query := `INSERT INTO links (id, created_at, short_link, long_link, owner_id) VALUES ($1, $2, $3, $4, $5)`

	_, err := l.db.Exec(query, link.ID, link.CreatedAt, link.ShortLink, link.LongLink, link.OwnerID)
	if err != nil {
		return fmt.Errorf("can't insert link: %v", err)
	}

	return nil
}

func (l *LinkStore) Select(query string) ([]link.Link, error) {
	links := make([]link.Link, 0)
	rows, err := l.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("can't select links: %v", err)
	}

	for rows.Next() {
		l := link.Link{}
		err := rows.Scan(&l.ID, &l.CreatedAt, &l.ShortLink, &l.LongLink, &l.OwnerID)
		if err != nil {
			return nil, fmt.Errorf("can't scan link: %v", err)
		}
		links = append(links, l)
	}

	return links, nil
}

func (l *LinkStore) Update(query string) error {
	_, err := l.db.Exec(query)
	if err != nil {
		return fmt.Errorf("can't update link: %v", err)
	}

	return nil
}
