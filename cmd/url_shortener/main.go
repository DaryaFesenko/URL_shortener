package main

import (
	"fmt"
	"log"
	"os"
	"urlshortener/app/starter"
	"urlshortener/db/pgsql"
)

func main() {
	os.Setenv("PG_DSN", "postgres://postgres:@localhost/shortener?sslmode=disable")
	app, err := starter.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	db, err := pgsql.InitDB(app.Config.Postgres)
	if err != nil {
		fmt.Println("db eroor")
		log.Fatal(err)
	}

	storers := pgsql.InitStorers(db)
	h := app.InitServices(storers)

	app.Serve(app.Config.ServerAddress, h)
}
