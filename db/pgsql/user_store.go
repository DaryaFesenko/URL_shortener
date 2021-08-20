package pgsql

import (
	"database/sql"
	"fmt"
	"urlshortener/app/services/user"
)

var _ user.UserStorer = &UserStore{}

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}

func (l *UserStore) Insert(user user.User) error {
	query := `INSERT INTO users (id, login, password) VALUES ($1, $2, $3)`

	_, err := l.db.Exec(query, user.ID, user.Login, user.Password)
	if err != nil {
		return fmt.Errorf("can't insert user: %v", err)
	}

	return nil
}

func (l *UserStore) Select(query string, params ...interface{}) ([]user.User, error) {
	users := make([]user.User, 0)
	rows, err := l.db.Query(query, params...)
	if err != nil {
		return nil, fmt.Errorf("can't select user: %v", err)
	}

	for rows.Next() {
		u := user.User{}
		err := rows.Scan(&u.ID, &u.Login, &u.Password)
		if err != nil {
			return nil, fmt.Errorf("can't scan user: %v", err)
		}
		users = append(users, u)
	}

	return users, nil
}

func (l *UserStore) Update(query string, params ...interface{}) error {
	_, err := l.db.Exec(query, params...)
	if err != nil {
		return fmt.Errorf("can't update user: %v", err)
	}

	return nil
}
