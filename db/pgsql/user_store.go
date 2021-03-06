package pgsql

import (
	"database/sql"
	"fmt"
	"urlshortener/app/services/auth"
)

var _ auth.UserStorer = &UserStore{}

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}

func (u *UserStore) ExistUserByLogin(login string) (bool, error) {
	users := make([]auth.User, 0)
	query := `SELECT login FROM users WHERE login = $1`
	rows, err := u.db.Query(query, login)
	if err != nil {
		return false, fmt.Errorf("can't select user: %v", err)
	}

	for rows.Next() {
		u := auth.User{}
		err := rows.Scan(&u.ID, &u.Login, &u.Password)
		if err != nil {
			return false, fmt.Errorf("can't scan user: %v", err)
		}
		users = append(users, u)
	}

	if len(users) == 0 {
		return false, nil
	}

	return true, nil
}

func (u *UserStore) Get(login, password string) (*auth.User, error) {
	users := make([]auth.User, 0)
	query := `SELECT id, login, password FROM users WHERE login = $1 AND password = $2`
	rows, err := u.db.Query(query, login, password)
	if err != nil {
		return nil, fmt.Errorf("can't select user: %v", err)
	}

	for rows.Next() {
		u := auth.User{}
		err := rows.Scan(&u.ID, &u.Login, &u.Password)
		if err != nil {
			return nil, fmt.Errorf("can't scan user: %v", err)
		}
		users = append(users, u)
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return &users[0], nil
}

func (l *UserStore) Insert(user *auth.User) error {
	query := `INSERT INTO users (id, login, password) VALUES ($1, $2, $3)`

	_, err := l.db.Exec(query, &user.ID, &user.Login, &user.Password)
	if err != nil {
		return fmt.Errorf("can't insert user: %v", err)
	}

	return nil
}

func (l *UserStore) Select(query string, params ...interface{}) ([]auth.User, error) {
	users := make([]auth.User, 0)
	rows, err := l.db.Query(query, params...)
	if err != nil {
		return nil, fmt.Errorf("can't select user: %v", err)
	}

	for rows.Next() {
		u := auth.User{}
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
