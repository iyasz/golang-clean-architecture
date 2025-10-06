package test

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/iyasz/golang-clean-architecture/internal/config"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var log *logrus.Logger
var db *gorm.DB
var validate *validator.Validate
var app *chi.Mux

func init() {
	conf := config.Load("test")
	log = config.NewLlogger(&conf.Logrus)
	db = config.NewDatabase(&conf.Database, log)
	validate = config.NewValidator()
	app = config.NewChi(conf)

	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		App:      app,
		Log:      log,
		Validate: validate,
	})
}
