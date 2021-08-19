package app

import (
	"database/sql"
	"urlshortener/pkg/link"
	"urlshortener/pkg/user"
)

type App struct {
	config *Config

	userService *user.UserService
	linkService *link.LinkService
	//linkTransitService *linktransition.LinkTransitService
}

func NewApp(configPath string) (*App, error) {
	a := App{}

	c, err := NewConfig(configPath)
	if err != nil {
		return &a, err
	}
	a.config = c

	db, err := a.initDB(c.GetPostgres())
	if err != nil {
		return &a, err
	}

	a.userService = user.NewUserService(db)
	a.linkService = link.NewLinkService(db)

	return &a, nil
}

func (a *App) initDB(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return db, err
	}

	return db, nil
}

func (a *App) Run() {
	users, err := a.userService.Store.Select("SELECT * FROM users")
	if err != nil {
		//logs
	}

	user := users[0]

	query := `UPDATE users SET login = $1 WHERE id = $2`
	err = a.userService.Store.Update(query, "test", user.ID)
	if err != nil {
		//logs
	}
}
