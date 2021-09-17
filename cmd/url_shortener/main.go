package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
	"urlshortener/api/server"
	"urlshortener/app/starter"
	"urlshortener/db/pgsql"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	go watchSignal(cancel)

	app, err := starter.NewApp(getConfigPath())
	if err != nil {
		log.Fatal(err)
	}

	db, err := pgsql.InitDB(app.Config.GetPostgres())
	if err != nil {
		log.Fatal(err)
	}

	storers := pgsql.InitStorers(db)
	h := app.InitServices(storers)
	s := server.NewServer(app.Config.ServerAddress, h)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go app.Serve(ctx, wg, s)

	<-ctx.Done()
	cancel()
	wg.Wait()
}

func getConfigPath() string {
	var configPath string

	flag.StringVar(&configPath, "c", "config.yaml", "Used for set path to config file.")
	flag.Parse()

	return configPath
}

func watchSignal(cancel context.CancelFunc) {
	osSignalChan := make(chan os.Signal, 1)
	signal.Notify(osSignalChan, os.Interrupt)

	<-osSignalChan

	log.Println("user interrupted")
	cancel()
}
