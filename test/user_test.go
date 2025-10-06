package test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/iyasz/golang-clean-architecture/internal/entity"
	"github.com/iyasz/golang-clean-architecture/internal/model"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestRegister(t *testing.T) {
	ClearAll()
	requestBody := model.RegisterUserRequest{
		ID:       "khannedy",
		Password: "rahasia",
		Name:     "Eko Khannedy",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	server := httptest.NewServer(app)
	defer server.Close()

	req, err := http.NewRequest(http.MethodPost, server.URL+BaseUsersAPIURL, strings.NewReader(string(bodyJson)))
	assert.Nil(t, err)
	SetupHeader(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, requestBody.ID, responseBody.Data.ID)
	assert.Equal(t, requestBody.Name, responseBody.Data.Name)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
}

func TestRegisterError(t *testing.T) {
	ClearAll()
	requestBody := model.RegisterUserRequest{
		ID:       "",
		Password: "",
		Name:     "",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	server := httptest.NewServer(app)
	defer server.Close()

	req, err := http.NewRequest(http.MethodPost, server.URL+BaseUsersAPIURL, strings.NewReader(string(bodyJson)))
	assert.Nil(t, err)
	SetupHeader(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.NotNil(t, responseBody.Errors)
}

func TestRegisterDuplicate(t *testing.T) {
	ClearAll()
	TestRegister(t)

	requestBody := model.RegisterUserRequest{
		ID:       "khannedy",
		Password: "rahasia",
		Name:     "Eko Khannedy",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	server := httptest.NewServer(app)
	defer server.Close()

	req, err := http.NewRequest(http.MethodPost, server.URL+BaseUsersAPIURL, strings.NewReader(string(bodyJson)))
	assert.Nil(t, err)
	SetupHeader(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusConflict, resp.StatusCode)
	assert.NotNil(t, responseBody.Errors)
}

func TestLogin(t *testing.T) {
	TestRegister(t)

	requestBody := model.LoginUserRequest{
		ID:       "khannedy",
		Password: "rahasia",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	server := httptest.NewServer(app)
	defer server.Close()

	req, err := http.NewRequest(http.MethodPost, server.URL+BaseUsersAPIURL+"/_login", strings.NewReader(string(bodyJson)))
	assert.Nil(t, err)
	SetupHeader(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, responseBody.Data.Token)

	user := new(entity.User)
	err = db.Where("id = ?", requestBody.ID).First(user).Error
	assert.Nil(t, err)
	assert.Equal(t, user.Token, responseBody.Data.Token)
}

func TestLoginWrongUsername(t *testing.T) {
	ClearAll()
	TestRegister(t)

	requestBody := model.LoginUserRequest{
		ID:       "wrong",
		Password: "rahasia",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	server := httptest.NewServer(app)
	defer server.Close()

	req, err := http.NewRequest(http.MethodPost, server.URL+BaseUsersAPIURL+"/_login", strings.NewReader(string(bodyJson)))
	assert.Nil(t, err)
	SetupHeader(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	assert.NotNil(t, responseBody.Errors)
}

func TestLoginWrongPassword(t *testing.T) {
	ClearAll()
	TestRegister(t)

	requestBody := model.LoginUserRequest{
		ID:       "khannedy",
		Password: "wrong",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	server := httptest.NewServer(app)
	defer server.Close()

	req, err := http.NewRequest(http.MethodPost, server.URL+BaseUsersAPIURL+"/_login", strings.NewReader(string(bodyJson)))
	assert.Nil(t, err)
	SetupHeader(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	assert.NotNil(t, responseBody.Errors)
}

func TestLogout(t *testing.T) {
	ClearAll()
	TestLogin(t) 

	user := new(entity.User)
	err := db.Where("id = ?", "khannedy").First(user).Error
	assert.Nil(t, err)

	server := httptest.NewServer(app)
	defer server.Close()

	req, err := http.NewRequest(http.MethodDelete, server.URL+BaseUsersAPIURL, nil)
	assert.Nil(t, err)
	SetupHeader(req)
	req.Header.Set("Authorization", user.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[bool])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.True(t, responseBody.Data)
}

func TestLogoutWrongAuthorization(t *testing.T) {
	ClearAll()
	TestLogin(t)

	server := httptest.NewServer(app)
	defer server.Close()

	req, err := http.NewRequest(http.MethodDelete, server.URL+BaseUsersAPIURL, nil)
	assert.Nil(t, err)
	SetupHeader(req)
	req.Header.Set("Authorization", "wrong")

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[bool])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	assert.NotNil(t, responseBody.Errors)
}

func TestGetCurrentUser(t *testing.T) {
	ClearAll()
	TestLogin(t)

	user := new(entity.User)
	err := db.Where("id = ?", "khannedy").First(user).Error
	assert.Nil(t, err)

	server := httptest.NewServer(app)
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL+BaseUsersAPIURL+"/_current", nil)
	assert.Nil(t, err)
	SetupHeader(req)
	req.Header.Set("Authorization", user.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, user.ID, responseBody.Data.ID)
	assert.Equal(t, user.Name, responseBody.Data.Name)
	assert.Equal(t, user.CreatedAt, responseBody.Data.CreatedAt)
	assert.Equal(t, user.UpdatedAt, responseBody.Data.UpdatedAt)
}

func TestGetCurrentUserFailed(t *testing.T) {
	ClearAll()
	TestLogin(t) 

	server := httptest.NewServer(app)
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL+BaseUsersAPIURL+"/_current", nil)
	assert.Nil(t, err)
	SetupHeader(req)
	req.Header.Set("Authorization", "wrong")

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	assert.NotNil(t, responseBody.Errors)
}

func TestUpdateUserName(t *testing.T) {
	ClearAll()
	TestLogin(t) 

	user := new(entity.User)
	err := db.Where("id = ?", "khannedy").First(user).Error
	assert.Nil(t, err)

	requestBody := model.UpdateUserRequest{
		Name: "Eko Kurniawan Khannedy",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	server := httptest.NewServer(app)
	defer server.Close()

	req, err := http.NewRequest(http.MethodPatch, server.URL+BaseUsersAPIURL+"/_current", strings.NewReader(string(bodyJson)))
	assert.Nil(t, err)
	SetupHeader(req)
	req.Header.Set("Authorization", user.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, user.ID, responseBody.Data.ID)
	assert.Equal(t, requestBody.Name, responseBody.Data.Name)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
}

func TestUpdateUserPassword(t *testing.T) {
	ClearAll()
	TestLogin(t) 

	user := new(entity.User)
	err := db.Where("id = ?", "khannedy").First(user).Error
	assert.Nil(t, err)

	requestBody := model.UpdateUserRequest{
		Password: "rahasialagi",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	server := httptest.NewServer(app)
	defer server.Close()

	req, err := http.NewRequest(http.MethodPatch, server.URL+BaseUsersAPIURL+"/_current", strings.NewReader(string(bodyJson)))
	assert.Nil(t, err)
	SetupHeader(req)
	req.Header.Set("Authorization", user.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, user.ID, responseBody.Data.ID)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)

	user = new(entity.User)
	err = db.Where("id = ?", "khannedy").First(user).Error
	assert.Nil(t, err)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password))
	assert.Nil(t, err)
}

func TestUpdateFailed(t *testing.T) {
	ClearAll()
	TestLogin(t) 

	requestBody := model.UpdateUserRequest{
		Password: "rahasialagi",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	server := httptest.NewServer(app)
	defer server.Close()

	req, err := http.NewRequest(http.MethodPatch, server.URL+BaseUsersAPIURL+"/_current", strings.NewReader(string(bodyJson)))
	assert.Nil(t, err)
	SetupHeader(req)
	req.Header.Set("Authorization", "wrong")

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	assert.NotNil(t, responseBody.Errors)
}