package config

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/iyasz/golang-clean-architecture/internal/delivery/http"
	"github.com/iyasz/golang-clean-architecture/internal/delivery/http/middleware"
	"github.com/iyasz/golang-clean-architecture/internal/repository"
	"github.com/iyasz/golang-clean-architecture/internal/usecase"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB  *gorm.DB
	App *chi.Mux
	Log *logrus.Logger
	Validate *validator.Validate
}

func Bootstrap(config *BootstrapConfig) {

	// setup repositories
	userRepository := repository.NewUserRepository(config.Log)
	contactRepository := repository.NewContactRepository(config.Log)
	addressRepository := repository.NewAddressRepository(config.Log)

	// setup use cases
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validate, userRepository)
	contactUseCase := usecase.NewContactUseCase(config.DB, config.Log, config.Validate, contactRepository)
	addressUseCase := usecase.NewAddressUseCase(config.DB, config.Log, config.Validate, addressRepository, contactRepository)

	// setup controller 
	userController := http.NewUserController(config.Log, userUseCase)

	// setup middleware
	authMiddleware := middleware.NewAuth(userUseCase)

}