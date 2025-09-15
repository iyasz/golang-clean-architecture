package main

import (
	"github.com/iyasz/golang-clean-architecture/internal/config"
)

func main() {

	conf := config.Load()

	// Logging 
	log := config.NewLlogger(&conf.Logrus)
	
	// Database 
	db := config.NewDatabase(&conf.Database, log)

	// Validator 
	validate := config.NewValidator()

	app := config.NewChi(conf)

	config.Bootstrap(&config.BootstrapConfig{
		DB: db,
		App: app,
		Log: log,
		Validate: validate,
	})

	log.Fatal(config.StartServer(app, &conf.Server, log))
}