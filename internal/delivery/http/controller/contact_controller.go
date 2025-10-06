package controller

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/iyasz/golang-clean-architecture/internal/delivery/http/middleware"
	"github.com/iyasz/golang-clean-architecture/internal/helper"
	"github.com/iyasz/golang-clean-architecture/internal/model"
	"github.com/iyasz/golang-clean-architecture/internal/usecase"
	"github.com/sirupsen/logrus"
)

type ContactController struct {
	Log *logrus.Logger
	ContactUseCase *usecase.ContactUseCase
}


func NewContactController(log *logrus.Logger, contactUseCase *usecase.ContactUseCase) *ContactController {
	return &ContactController{
		Log: log,
		ContactUseCase: contactUseCase,
	}
}

func (c *ContactController) Create(w http.ResponseWriter, r *http.Request) {
	auth := middleware.GetUser(r)

	request := new(model.CreateContactRequest)
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		helper.ErrorResponse(w, helper.ErrBadRequest)
		return
	}
	request.UserId = auth.ID

	response, err := c.ContactUseCase.Create(r.Context(), request)
	if err != nil {
		c.Log.WithError(err).Error("error creating contact")
		helper.ErrorResponse(w, err)
		return
	}

	helper.SuccessResponse(w, model.WebResponse[*model.ContactResponse]{Data: response}, http.StatusCreated)
}

func (c *ContactController) List(w http.ResponseWriter, r *http.Request) {
	auth := middleware.GetUser(r)

	name := r.URL.Query().Get("name")
	email := r.URL.Query().Get("email")
	phone := r.URL.Query().Get("phone")

	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			page = parsed
		}
	}

	size := 10
	if s := r.URL.Query().Get("size"); s != "" {
		if parsed, err := strconv.Atoi(s); err == nil {
			size = parsed
		}
	}

	request := &model.SearchContactRequest{
		UserId: auth.ID,
		Name:   name,
		Email: 	email,
		Phone:  phone,
		Page:   page,
		Size:   size,
	}

	responses, total, err := c.ContactUseCase.Search(r.Context(), request)
	if err != nil {
		c.Log.WithError(err).Error("error searching contact")
		helper.ErrorResponse(w, err)
		return
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	helper.SuccessResponse(w, model.WebResponse[[]model.ContactResponse]{
		Data: responses,
		Paging: paging,
		}, http.StatusOK)

}

func (c *ContactController) Get(w http.ResponseWriter, r *http.Request) {
	auth := middleware.GetUser(r)

	request := &model.GetContactRequest{
		UserId: auth.ID,
		ID:     chi.URLParam(r, "contactId"),
	}

	response, err := c.ContactUseCase.Get(r.Context(), request)
	if err != nil {
		c.Log.WithError(err).Error("error getting contact")
		helper.ErrorResponse(w, err)
		return
	}

	helper.SuccessResponse(w, model.WebResponse[*model.ContactResponse]{Data: response}, http.StatusOK)
}

func (c *ContactController) Update(w http.ResponseWriter, r *http.Request) {
	auth := middleware.GetUser(r)

	request := new(model.UpdateContactRequest)
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		helper.ErrorResponse(w, helper.ErrBadRequest)
		return
	}

	request.UserId = auth.ID
	request.ID = chi.URLParam(r, "contactId")

	response, err := c.ContactUseCase.Update(r.Context(), request)
	if err != nil {
		c.Log.WithError(err).Error("error updating contact")
		helper.ErrorResponse(w, err)
		return
	}

	helper.SuccessResponse(w, model.WebResponse[*model.ContactResponse]{Data: response}, http.StatusOK)
}

func (c *ContactController) Delete(w http.ResponseWriter, r *http.Request) {
	auth := middleware.GetUser(r)
	contactId := chi.URLParam(r, "contactId")

	request := &model.DeleteContactRequest{
		UserId: auth.ID,
		ID:     contactId,
	}

	if err := c.ContactUseCase.Delete(r.Context(), request); err != nil {
		c.Log.WithError(err).Error("error deleting contact")
		helper.ErrorResponse(w, err)
		return
	}

	helper.SuccessResponse(w, model.WebResponse[bool]{Data: true}, http.StatusOK)
}
