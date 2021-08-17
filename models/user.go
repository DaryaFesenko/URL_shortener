package models

import uuid "github.com/satori/go.uuid"

type User struct {
	ID       uuid.UUID
	Login    string
	Password string
}

type UserStore struct {
}

func (l *UserStore) Insert(user User) error {
	query := `INSERT INTO users (id, login, password) VALUES ($1, $2, $3)`

	_, err := DB.Exec(query, user.ID, user.Login, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (l *UserStore) Select(query string, params ...interface{}) ([]User, error) {
	users := make([]User, 0)
	rows, err := DB.Query(query, params...)
	if err != nil {
		return users, err
	}

	for rows.Next() {
		u := User{}
		err := rows.Scan(&u.ID, &u.Login, &u.Password)
		if err != nil {
			return users, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (l *UserStore) Update(query string, params ...interface{}) error {
	_, err := DB.Exec(query, params...)
	if err != nil {
		return err
	}

	return nil
}
