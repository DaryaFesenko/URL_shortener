package linktransition

import (
	"database/sql"
	"urlshortener/models"
)

type LinkTransitionStore struct {
	db *sql.DB
}

func NewLinkTransitStorer(db *sql.DB) *LinkTransitionStore {
	return &LinkTransitionStore{db: db}
}

func (l *LinkTransitionStore) Insert(lt models.LinkTransition) error {
	query := `INSERT INTO link_transitions (id, link_id, used_user_id, used_count) VALUES ($1, $2, $3, $4)`

	_, err := l.db.Exec(query, lt.ID, lt.LinkID, lt.UsedUserID, lt.UsedCount)
	if err != nil {
		return err
	}

	return nil
}

func (l *LinkTransitionStore) Select(query string) ([]models.LinkTransition, error) {
	lt := make([]models.LinkTransition, 0)
	rows, err := l.db.Query(query)
	if err != nil {
		return lt, err
	}

	for rows.Next() {
		l := models.LinkTransition{}
		err := rows.Scan(&l.ID, &l.LinkID, &l.UsedUserID, &l.UsedCount)
		if err != nil {
			return lt, err
		}
		lt = append(lt, l)
	}

	return lt, nil
}

func (l *LinkTransitionStore) Update(query string) error {
	_, err := l.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
