package config

import "github.com/sirupsen/logrus"

func NewLlogger(cfl *Logrus) *logrus.Logger {
	log := logrus.New()

	log.SetLevel(logrus.Level(cfl.Level))
	log.SetFormatter(&logrus.JSONFormatter{})

	return log
}