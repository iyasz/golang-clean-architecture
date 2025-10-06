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


	c.App.Post("/api/users", c.UserController.Register)
	c.App.Post("/api/users/_login", c.UserController.Login)

	c.App.Route("/api", func(r chi.Router) {
		r.Use(c.AuthMiddleware)
		r.Delete("/users", c.UserController.Logout)
		r.Patch("/users/_current", c.UserController.Update)
		r.Get("/users/_current", c.UserController.Current)

		r.Get("/contacts", c.ContactController.List)
		r.Post("/contacts", c.ContactController.Create)
		r.Put("/contacts/{contactId}", c.ContactController.Update)
		r.Get("/contacts/{contactId}", c.ContactController.Get)
		r.Delete("/contacts/{contactId}", c.ContactController.Delete)

		r.Get("/contacts/{contactId}/addresses", c.AddressController.List)
		r.Post("/contacts/{contactId}/addresses", c.AddressController.Create)
		r.Put("/contacts/{contactId}/addresses/{addressId}", c.AddressController.Update)
		r.Get("/contacts/{contactId}/addresses/{addressId}", c.AddressController.Get)
		r.Delete("/contacts/{contactId}/addresses/{addressId}", c.AddressController.Delete)
	})
}

