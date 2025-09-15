package http

import (
	"encoding/json"
	"net/http"

	"github.com/iyasz/golang-clean-architecture/internal/delivery/http/middleware"
	"github.com/iyasz/golang-clean-architecture/internal/helper"
	"github.com/iyasz/golang-clean-architecture/internal/model"
	"github.com/iyasz/golang-clean-architecture/internal/usecase"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	Log         *logrus.Logger
	UserUseCase *usecase.UserUseCase
}

func NewUserController(log *logrus.Logger, userUseCase *usecase.UserUseCase) *UserController {
	return &UserController{
		Log:         log,
		UserUseCase: userUseCase,
	}
}

func (c *UserController) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	request := new(model.RegisterUserRequest)
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		helper.ErrorResponse(w, helper.ErrBadRequest)
		return
	}

	response, err := c.UserUseCase.Create(ctx, request)
	if err != nil {
		c.Log.Warnf("Failed to register user : %+v", err)
		helper.ErrorResponse(w, err)
		return
	}

	helper.SuccessResponse(w, model.WebResponse[*model.UserResponse]{Data: response}, http.StatusCreated)
}

func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	request := new(model.LoginUserRequest)
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		helper.ErrorResponse(w, helper.ErrBadRequest)
		return
	}

	response, err := c.UserUseCase.Login(ctx, request)
	if err != nil {
		c.Log.Warnf("Failed to login user : %+v", err)
		helper.ErrorResponse(w, err)
		return
	}

	helper.SuccessResponse(w, model.WebResponse[*model.UserResponse]{Data: response}, http.StatusOK)
}

func (c *UserController) Current(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	auth := middleware.GetUser(r)

	request := &model.GetUserRequest{
		ID: auth.ID,
	}

	response, err := c.UserUseCase.Current(ctx, request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to get current user")
		helper.ErrorResponse(w, err)
		return
	}

	helper.SuccessResponse(w, model.WebResponse[*model.UserResponse]{Data: response}, http.StatusOK)
}

func (c *UserController) Logout(w http.ResponseWriter, r *http.Request) {
	auth := middleware.GetUser(r)

	request := &model.LogoutUserRequest{
		ID: auth.ID,
	}

	response, err := c.UserUseCase.Logout(r.Context(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to logout user")
		helper.ErrorResponse(w, err)
		return
	}

	helper.SuccessResponse(w, model.WebResponse[bool]{Data: response}, http.StatusOK)
}

func (c *UserController) Update(w http.ResponseWriter, r *http.Request) {
	auth := middleware.GetUser(r)

	request := new(model.UpdateUserRequest)

	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		helper.ErrorResponse(w, helper.ErrBadRequest)
		return
	}

	request.ID = auth.ID
	response, err := c.UserUseCase.Update(r.Context(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to update user")
		helper.ErrorResponse(w, err)
		return
	}

	helper.SuccessResponse(w, model.WebResponse[*model.UserResponse]{Data: response}, http.StatusOK)
}
