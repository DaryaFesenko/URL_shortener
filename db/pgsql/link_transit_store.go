package pgsql

import (
	"database/sql"
	"fmt"
	"time"
	"urlshortener/app/services/link"

	"github.com/google/uuid"
)

var _ link.LinkTransitStorer = &LinkTransitionStore{}

type LinkTransitionStore struct {
	db *sql.DB
}

func NewLinkTransitStorer(db *sql.DB) *LinkTransitionStore {
	return &LinkTransitionStore{db: db}
}

func (l *LinkTransitionStore) Insert(lt link.LinkTransition) error {
	query := `INSERT INTO link_transitions (id, link_id, ip, used_count, date) VALUES ($1, $2, $3, $4, $5)`

	_, err := l.db.Exec(query, lt.ID, lt.LinkID, lt.IP, lt.UsedCount, lt.Date)
	if err != nil {
		return fmt.Errorf("can't insert link transit: %v", err)
	}

	return nil
}

func (l *LinkTransitionStore) UpdateTransitCount(id uuid.UUID, usedCount int) error {
	_, err := l.db.Exec(`UPDATE link_transitions SET used_count = $2, date = $3 WHERE id = $1`, id, usedCount, time.Now())
	if err != nil {
		return fmt.Errorf("can't update link transit: %v", err)
	}

	return nil
}

func (l *LinkTransitionStore) StatisticLink(linkID uuid.UUID) ([]link.LinkTransition, error) {
	linkTransitions, err := l.getTransitions(`SELECT id, link_id, ip, used_count, date FROM link_transitions WHERE link_id = $1`, linkID)
	if err != nil {
		return nil, err
	}

	return linkTransitions, nil
}

func (l *LinkTransitionStore) GetTransit(usedUserID string, linkID uuid.UUID) (link.LinkTransition, error) {
	linkTransitions, err := l.getTransitions(
		`SELECT id, link_id, ip, used_count, date FROM link_transitions WHERE link_id = $1 AND ip = $2`,
		linkID, usedUserID)

	if err != nil {
		return link.LinkTransition{}, err
	}

	if len(linkTransitions) == 0 {
		return link.LinkTransition{}, sql.ErrNoRows
	}

	return linkTransitions[0], nil
}

func (l *LinkTransitionStore) DeleteLinkTransit(linkID uuid.UUID) error {
	_, err := l.db.Exec(`DELETE FROM link_transitions WHERE link_id = $1`, linkID)
	if err != nil {
		return fmt.Errorf("can't delete link transitions: %v", err)
	}

	return nil
}

// nolint:dupl // it's NOT duplicate
func (l *LinkTransitionStore) getTransitions(query string, params ...interface{}) ([]link.LinkTransition, error) {
	lt := make([]link.LinkTransition, 0)
	rows, err := l.db.Query(query, params...)
	if err != nil {
		return nil, fmt.Errorf("can't select link transit: %v", err)
	}

	for rows.Next() {
		l := link.LinkTransition{}
		err := rows.Scan(&l.ID, &l.LinkID, &l.IP, &l.UsedCount, &l.Date)
		if err != nil {
			return nil, fmt.Errorf("can't scan link transit: %v", err)
		}
		lt = append(lt, l)
	}

	return lt, nil
}
