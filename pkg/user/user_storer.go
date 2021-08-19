package user

import (
	"database/sql"
	"urlshortener/models"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}

func (l *UserStore) Insert(user models.User) error {
	query := `INSERT INTO users (id, login, password) VALUES ($1, $2, $3)`

	_, err := l.db.Exec(query, user.ID, user.Login, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (l *UserStore) Select(query string, params ...interface{}) ([]models.User, error) {
	users := make([]models.User, 0)
	rows, err := l.db.Query(query, params...)
	if err != nil {
		return users, err
	}

	for rows.Next() {
		u := models.User{}
		err := rows.Scan(&u.ID, &u.Login, &u.Password)
		if err != nil {
			return users, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (l *UserStore) Update(query string, params ...interface{}) error {
	_, err := l.db.Exec(query, params...)
	if err != nil {
		return err
	}

	return nil
}
