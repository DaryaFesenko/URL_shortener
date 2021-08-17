package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Link struct {
	ID        uuid.UUID
	CreatedAt time.Time

	ShortLink string
	LongLink  string

	OwnerID uuid.UUID
}

type LinkStore struct {
}

func (l *LinkStore) Insert(link Link) error {
	query := `INSERT INTO links (id, created_at, short_link, long_link, owner_id) VALUES ($1, $2, $3, $4, $5)`

	_, err := DB.Exec(query, link.ID, link.CreatedAt, link.ShortLink, link.LongLink, link.OwnerID)
	if err != nil {
		return err
	}

	return nil
}

func (l *LinkStore) Select(query string) ([]Link, error) {
	links := make([]Link, 0)
	rows, err := DB.Query(query)
	if err != nil {
		return links, err
	}

	for rows.Next() {
		l := Link{}
		err := rows.Scan(&l.ID, &l.CreatedAt, &l.ShortLink, &l.LongLink, &l.OwnerID)
		if err != nil {
			return links, err
		}
		links = append(links, l)
	}

	return links, nil
}

func (l *LinkStore) Update(query string) error {
	_, err := DB.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
