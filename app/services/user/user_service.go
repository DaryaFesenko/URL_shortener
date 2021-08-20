package user

type UserStorer interface {
	Insert(user User) error
	Select(query string, params ...interface{}) ([]User, error)
	Update(query string, params ...interface{}) error
}

type UserService struct {
	Store UserStorer
}

func NewUserService(store UserStorer) *UserService {
	return &UserService{Store: store}
}
