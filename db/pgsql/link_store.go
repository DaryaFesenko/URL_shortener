package pgsql

import (
	"database/sql"
	"fmt"
	"urlshortener/app/services/link"
	linkService "urlshortener/app/services/link"

	"github.com/google/uuid"
)

var _ linkService.LinkStorer = &LinkStore{}

type LinkStore struct {
	db *sql.DB
}

func NewLinkStorer(db *sql.DB) *LinkStore {
	return &LinkStore{db: db}
}

func (l *LinkStore) Insert(link *linkService.Link) error {
	query := `INSERT INTO links (id, created_at, short_link, long_link, owner_id) VALUES ($1, $2, $3, $4, $5)`

	_, err := l.db.Exec(query, link.ID, link.CreatedAt, link.ShortLink, link.LongLink, link.OwnerID)
	if err != nil {
		return fmt.Errorf("can't insert link: %v", err)
	}

	return nil
}

func (l *LinkStore) ExistLongLink(userID *uuid.UUID, longLink string) (bool, error) {
	links := make([]linkService.Link, 0)
	query := `SELECT * FROM links WHERE owner_id = $1 AND long_link = $2`
	rows, err := l.db.Query(query, &userID, &longLink)
	if err != nil {
		return false, fmt.Errorf("can't select link: %v", err)
	}

	for rows.Next() {
		link := linkService.Link{}
		err := rows.Scan(&link.ID, &link.CreatedAt, &link.ShortLink, &link.LongLink, &link.OwnerID)
		if err != nil {
			return false, fmt.Errorf("can't scan link: %v", err)
		}
		links = append(links, link)
	}

	if len(links) != 0 {
		return true, nil
	}

	return false, nil
}

// nolint:dupl // it's NOT duplicate
func (l *LinkStore) getLinks(query string, params ...interface{}) ([]link.Link, error) {
	links := make([]link.Link, 0)
	rows, err := l.db.Query(query, params...)
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

func (l *LinkStore) GetLink(linkID uuid.UUID) (*link.Link, error) {
	links, err := l.getLinks(`SELECT id, created_at, short_link, long_link, owner_id FROM links WHERE id = $1`, &linkID)
	if err != nil {
		return nil, err
	}

	if len(links) == 0 {
		return nil, sql.ErrNoRows
	}

	return &links[0], nil
}
func (l *LinkStore) GetLinks(ownerID uuid.UUID) ([]link.Link, error) {
	return l.getLinks(`SELECT id, created_at, short_link, long_link, owner_id FROM links WHERE owner_id = $1`, &ownerID)
}

func (l *LinkStore) GetLongLink(shortLink string) (string, error) {
	links, err := l.getLinks(`SELECT id, created_at, short_link, long_link, owner_id FROM links WHERE short_link = $1`, shortLink)
	if err != nil {
		return "", err
	}

	if len(links) == 0 {
		return "", fmt.Errorf("can't get long link with short link %s", shortLink)
	}

	return links[0].LongLink, nil
}

func (l *LinkStore) GetLinkIDByShortLink(shortLink string) (*uuid.UUID, error) {
	links, err := l.getLinks(`SELECT id, created_at, short_link, long_link, owner_id FROM links WHERE short_link = $1`, shortLink)
	if err != nil {
		return nil, err
	}

	if len(links) == 0 {
		return nil, fmt.Errorf("can't get link id with short link: %s", shortLink)
	}

	return &links[0].ID, nil
}

func (l *LinkStore) DeleteLink(linkID uuid.UUID) error {
	_, err := l.db.Exec(`DELETE FROM links WHERE id = $1`, linkID)
	if err != nil {
		return fmt.Errorf("can't delete link: %v", err)
	}

	return nil
}
