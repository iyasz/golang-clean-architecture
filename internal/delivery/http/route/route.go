package route

import (
	"net/http"

	"github.com/iyasz/golang-clean-architecture/internal/delivery/http/controller"
	"github.com/go-chi/chi/v5"
)

type RouteConfig struct {
	App               *chi.Mux
	UserController    *controller.UserController
	ContactController *controller.ContactController
	AddressController *controller.AddressController
	AuthMiddleware    func(http.Handler) http.Handler
}

func (c *RouteConfig) Setup() {

}

func (c *RouteConfig) SetupGuestRoute() {
	c.App.Post("/register", c.UserController.Register)
}
