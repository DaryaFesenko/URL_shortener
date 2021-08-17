package models

import uuid "github.com/satori/go.uuid"

type LinkTransition struct {
	ID         uuid.UUID
	LinkID     uuid.UUID
	UsedUserID uuid.UUID
	UsedCount  int
}

type LinkTransitionStore struct {
}

func (l *LinkTransitionStore) Insert(lt LinkTransition) error {
	query := `INSERT INTO link_transitions (id, link_id, used_user_id, used_count) VALUES ($1, $2, $3, $4)`

	_, err := DB.Exec(query, lt.ID, lt.LinkID, lt.UsedUserID, lt.UsedCount)
	if err != nil {
		return err
	}

	return nil
}

func (l *LinkTransitionStore) Select(query string) ([]LinkTransition, error) {
	lt := make([]LinkTransition, 0)
	rows, err := DB.Query(query)
	if err != nil {
		return lt, err
	}

	for rows.Next() {
		l := LinkTransition{}
		err := rows.Scan(&l.ID, &l.LinkID, &l.UsedUserID, &l.UsedCount)
		if err != nil {
			return lt, err
		}
		lt = append(lt, l)
	}

	return lt, nil
}

func (l *LinkTransitionStore) Update(query string) error {
	_, err := DB.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
