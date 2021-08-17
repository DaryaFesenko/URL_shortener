package app

import (
	"urlshortener/models"
)

type App struct {
	config *Config
}

func NewApp(configPath string) (*App, error) {
	a := App{}

	c, err := NewConfig(configPath)
	if err != nil {
		return &a, err
	}
	a.config = c

	err = a.initDB()
	if err != nil {
		return &a, err
	}

	return &a, nil
}

func (a *App) initDB() error {
	err := models.InitDB(a.config.GetPostgres())

	if err != nil {
		return err
	}

	return nil
}

func (a *App) Run() {
	us := models.UserStore{}

	users, err := us.Select("SELECT * FROM users")
	if err != nil {
		//logs
	}

	user := users[0]

	query := `UPDATE users SET login = $1 WHERE id = $2`
	err = us.Update(query, "test", user.ID)
	if err != nil {
		//logs
	}
}
