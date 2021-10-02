package starter

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"urlshortener/api/handler"
	"urlshortener/app/services/auth"
	"urlshortener/app/services/link"
)

const (
	timeout = 60
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
	linkTransitStorer link.LinkTransitStorer
}

func NewStorers(user auth.UserStorer, link link.LinkStorer, lt link.LinkTransitStorer) *Storers {
	return &Storers{
		userStorer:        user,
		linkStorer:        link,
		linkTransitStorer: lt,
	}
}

func NewApp() (*App, error) {
	a := App{}

	c, err := NewConfig()
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
	linkService := link.NewLinkService(a.storers.linkStorer, a.storers.linkTransitStorer)

	return handler.NewRouter(authService, linkService, sign_key, a.Config.ServerAddress)
}

func (a *App) Serve(addr string, h http.Handler) {
	srv := http.Server{
		Addr:              addr,
		Handler:           h,
		ReadTimeout:       timeout * time.Second,
		WriteTimeout:      timeout * time.Second,
		ReadHeaderTimeout: timeout * time.Second,
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		fmt.Println("start listen")

		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			if err == http.ErrServerClosed {
				log.Println(err.Error())
			} else {
				log.Println(err)
			}
		}
	}()

	go func() {
		defer wg.Done()

		signalListener := make(chan os.Signal, 4)
		signal.Notify(signalListener,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT)

		stop := <-signalListener
		log.Println("Received ", stop)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Println(err)
		}
	}()

	wg.Wait()
}
