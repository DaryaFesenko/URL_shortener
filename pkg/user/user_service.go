package user

import "database/sql"

type UserService struct {
	Store *UserStore
}

func NewUserService(db *sql.DB) *UserService {
	store := NewUserStore(db)

	return &UserService{Store: store}
}
