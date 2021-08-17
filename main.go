package main

import (
	"flag"
	"fmt"
	"urlshortener/app"
)

func main() {
	fmt.Println("starting course project")

	app, err := app.NewApp(getConfigPath())
	if err != nil {
		//logs
	}

	app.Run()
}

func getConfigPath() string {
	var configPath string

	flag.StringVar(&configPath, "c", "config.yaml", "Used for set path to config file.")
	flag.Parse()

	return configPath
}
