package link

import (
	"database/sql"
	"urlshortener/models"
)

type LinkStore struct {
	db *sql.DB
}

func NewLinkStorer(db *sql.DB) *LinkStore {
	return &LinkStore{db: db}
}

func (l *LinkStore) Insert(link models.Link) error {
	query := `INSERT INTO links (id, created_at, short_link, long_link, owner_id) VALUES ($1, $2, $3, $4, $5)`

	_, err := l.db.Exec(query, link.ID, link.CreatedAt, link.ShortLink, link.LongLink, link.OwnerID)
	if err != nil {
		return err
	}

	return nil
}

func (l *LinkStore) Select(query string) ([]models.Link, error) {
	links := make([]models.Link, 0)
	rows, err := l.db.Query(query)
	if err != nil {
		return links, err
	}

	for rows.Next() {
		l := models.Link{}
		err := rows.Scan(&l.ID, &l.CreatedAt, &l.ShortLink, &l.LongLink, &l.OwnerID)
		if err != nil {
			return links, err
		}
		links = append(links, l)
	}

	return links, nil
}

func (l *LinkStore) Update(query string) error {
	_, err := l.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
