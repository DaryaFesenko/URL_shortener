package starter

import (
	"context"
	"sync"
	"urlshortener/api/handler"
	"urlshortener/app/services/auth"
	"urlshortener/app/services/link"
	"urlshortener/app/services/linktransit"
)

type HTTPServer interface {
	Start(wg *sync.WaitGroup)
	Stop()
}
type App struct {
	Config *Config

	storers *Storers
}

type Storers struct {
	userStorer        auth.UserStorer
	linkStorer        link.LinkStorer
	linkTransitStorer linktransit.LinkTransitStorer
}

func NewStorers(user auth.UserStorer, link link.LinkStorer, lt linktransit.LinkTransitStorer) *Storers {
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

func (a *App) InitServices(s *Storers) *handler.Router {
	a.storers = s
	sign_key := make([]byte, 0)
	copy(sign_key, []byte(a.Config.SigningKey))

	authService := auth.NewAuthorizer(a.storers.userStorer, a.Config.HashSalt, sign_key, a.Config.ExpireDuration)
	linkService := link.NewLinkService(a.storers.linkStorer)
	//linkTransitService := linktransit.NewLinkTransitService(a.storers.linkTransitStorer)

	return handler.NewRouter(authService, linkService, sign_key)
}

func (a *App) Serve(ctx context.Context, wg *sync.WaitGroup, hs HTTPServer) {
	hs.Start(wg)
	<-ctx.Done()
	hs.Stop()
}
