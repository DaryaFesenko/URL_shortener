package main

import (
	"flag"
	"log"
	"urlshortener/app/starter"
	"urlshortener/db/pgsql"
)

func main() {
	app, err := starter.NewApp(getConfigPath())
	if err != nil {
		log.Fatal(err)
	}

	db, err := pgsql.InitDB(app.Config.GetPostgres())
	if err != nil {
		log.Fatal(err)
	}

	storers := pgsql.InitStorers(db)
	app.InitServices(storers)

	err = app.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func getConfigPath() string {
	var configPath string

	flag.StringVar(&configPath, "c", "config.yaml", "Used for set path to config file.")
	flag.Parse()

	return configPath
}
