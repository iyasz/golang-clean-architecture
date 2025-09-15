package config

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

func NewChi(conf *Config) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(60 * time.Second))

	// r.Use(middleware.SetHeader("Content-Type", "application/json"))

	return r
}

func StartServer(router *chi.Mux, config *Server, log *logrus.Logger) error {
	server := &http.Server{
		Addr:         config.Host + ":" + config.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Infof("Server starting on %s", server.Addr)
	return server.ListenAndServe()
}