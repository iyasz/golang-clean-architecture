package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/iyasz/golang-clean-architecture/internal/delivery/http/middleware"
	"github.com/iyasz/golang-clean-architecture/internal/helper"
	"github.com/iyasz/golang-clean-architecture/internal/model"
	"github.com/iyasz/golang-clean-architecture/internal/usecase"
	"github.com/sirupsen/logrus"
)

type AddressController struct {
	Log *logrus.Logger
	AddressUseCase *usecase.AddressUseCase
}

func NewAddressController(log *logrus.Logger, addressUseCase *usecase.AddressUseCase) *AddressController {
	return &AddressController{
		Log: log,
		AddressUseCase: addressUseCase,
	}
}

func (c *AddressController) Create(w http.ResponseWriter, r *http.Request) {
	auth := middleware.GetUser(r)

	request := new(model.CreateAddressRequest)
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		helper.ErrorResponse(w, helper.ErrBadRequest)
		return
	}

	request.UserId = auth.ID
	request.ContactId = chi.URLParam(r, "contactId")

	response, err := c.AddressUseCase.Create(r.Context(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to create address")
		helper.ErrorResponse(w, err)
		return
	}

	helper.SuccessResponse(w, model.WebResponse[*model.AddressResponse]{Data: response}, http.StatusCreated)
}

func (c *AddressController) List(w http.ResponseWriter, r *http.Request) {
	auth := middleware.GetUser(r)
	contactId := chi.URLParam(r, "contactId")

	request := &model.ListAddressRequest{
		UserId:    auth.ID,
		ContactId: contactId,
	}

	responses, err := c.AddressUseCase.List(r.Context(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to list addresses")
		helper.ErrorResponse(w, err)
		return
	}

	helper.SuccessResponse(w, model.WebResponse[[]model.AddressResponse]{Data: responses}, http.StatusOK)
}

func (c *AddressController) Get(w http.ResponseWriter, r *http.Request) {
	auth := middleware.GetUser(r)
	contactId := chi.URLParam(r, "contactId")
	addressId := chi.URLParam(r, "addressId")

	request := &model.GetAddressRequest{
		UserId:    auth.ID,
		ContactId: contactId,
		ID:        addressId,
	}

	response, err := c.AddressUseCase.Get(r.Context(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to get address")
		helper.ErrorResponse(w, err)
		return
	}

	helper.SuccessResponse(w, model.WebResponse[*model.AddressResponse]{Data: response}, http.StatusOK)
}

func (c *AddressController) Update(w http.ResponseWriter, r *http.Request) {
	auth := middleware.GetUser(r)

	request := new(model.UpdateAddressRequest)
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		helper.ErrorResponse(w, helper.ErrBadRequest)
		return
	}

	request.UserId = auth.ID
	request.ContactId = chi.URLParam(r, "contactId")
	request.ID = chi.URLParam(r, "addressId")

	response, err := c.AddressUseCase.Update(r.Context(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to update address")
		helper.ErrorResponse(w, err)
		return
	}

	helper.SuccessResponse(w, model.WebResponse[*model.AddressResponse]{Data: response}, http.StatusOK)
}

func (c *AddressController) Delete(w http.ResponseWriter, r *http.Request) {
	auth := middleware.GetUser(r)
	contactId := chi.URLParam(r, "contactId")
	addressId := chi.URLParam(r, "addressId")

	request := &model.DeleteAddressRequest{
		UserId:    auth.ID,
		ContactId: contactId,
		ID:        addressId,
	}

	if err := c.AddressUseCase.Delete(r.Context(), request); err != nil {
		c.Log.WithError(err).Error("failed to delete address")
		helper.ErrorResponse(w, err)
		return
	}

	helper.SuccessResponse(w, model.WebResponse[bool]{Data: true}, http.StatusNoContent)
}
