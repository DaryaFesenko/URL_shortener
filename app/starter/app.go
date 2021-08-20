package starter

import (
	"urlshortener/app/services/link"
	"urlshortener/app/services/linktransit"
	"urlshortener/app/services/user"
)

type App struct {
	Config *Config

	userService        *user.UserService
	linkService        *link.LinkService
	linkTransitService *linktransit.LinkTransitService

	storers *Storers
}

type Storers struct {
	userStorer        user.UserStorer
	linkStorer        link.LinkStorer
	linkTransitStorer linktransit.LinkTransitStorer
}

func NewStorers(user user.UserStorer, link link.LinkStorer, lt linktransit.LinkTransitStorer) *Storers {
	return &Storers{
		userStorer:        user,
		linkStorer:        link,
		linkTransitStorer: lt,
	}
}

func NewApp(configPath string) (*App, error) {
	a := App{}

	c, err := NewConfig(configPath)
	if err != nil {
		return &a, err
	}
	a.Config = c

	return &a, nil
}

func (a *App) InitServices(s *Storers) {
	a.storers = s

	a.userService = user.NewUserService(a.storers.userStorer)
	a.linkService = link.NewLinkService(a.storers.linkStorer)
	a.linkTransitService = linktransit.NewLinkTransitService(a.storers.linkTransitStorer)
}

func (a *App) Run() error {
	users, err := a.userService.Store.Select("SELECT * FROM users")
	if err != nil {
		return err
	}

	user := users[0]

	query := `UPDATE users SET login = $1 WHERE id = $2`
	err = a.userService.Store.Update(query, "test", user.ID)
	if err != nil {
		return err
	}

	return nil
}
