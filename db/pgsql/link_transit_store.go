package pgsql

import (
	"database/sql"
	"fmt"
	"urlshortener/app/services/linktransit"
)

var _ linktransit.LinkTransitStorer = &LinkTransitionStore{}

type LinkTransitionStore struct {
	db *sql.DB
}

func NewLinkTransitStorer(db *sql.DB) *LinkTransitionStore {
	return &LinkTransitionStore{db: db}
}

func (l *LinkTransitionStore) Insert(lt linktransit.LinkTransition) error {
	query := `INSERT INTO link_transitions (id, link_id, used_user_id, used_count) VALUES ($1, $2, $3, $4)`

	_, err := l.db.Exec(query, lt.ID, lt.LinkID, lt.UsedUserID, lt.UsedCount)
	if err != nil {
		return fmt.Errorf("can't insert link transit: %v", err)
	}

	return nil
}

func (l *LinkTransitionStore) Select(query string) ([]linktransit.LinkTransition, error) {
	lt := make([]linktransit.LinkTransition, 0)
	rows, err := l.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("can't select link transit: %v", err)
	}

	for rows.Next() {
		l := linktransit.LinkTransition{}
		err := rows.Scan(&l.ID, &l.LinkID, &l.UsedUserID, &l.UsedCount)
		if err != nil {
			return nil, fmt.Errorf("can't scan link transit: %v", err)
		}
		lt = append(lt, l)
	}

	return lt, nil
}

func (l *LinkTransitionStore) Update(query string) error {
	_, err := l.db.Exec(query)
	if err != nil {
		return fmt.Errorf("can't update link transit: %v", err)
	}

	return nil
}
