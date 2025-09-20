package test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/iyasz/golang-clean-architecture/internal/entity"
	"github.com/iyasz/golang-clean-architecture/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestCreateContact(t *testing.T) {
	TestLogin(t)

	user := new(entity.User)
	err := db.Where("id = ?", "khannedy").First(user).Error
	assert.Nil(t, err)

	requestBody := model.CreateContactRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "eko@example.com",
		Phone:     "088888888888",
	}
	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	// Buat server Chi untuk testing
	server := httptest.NewServer(app)
	defer server.Close()

	// Kirim request ke server
	req, err := http.NewRequest(http.MethodPost, server.URL+"/api/contacts", strings.NewReader(string(bodyJson)))
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", user.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.ContactResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, requestBody.FirstName, responseBody.Data.FirstName)
	assert.Equal(t, requestBody.LastName, responseBody.Data.LastName)
	assert.Equal(t, requestBody.Email, responseBody.Data.Email)
	assert.Equal(t, requestBody.Phone, responseBody.Data.Phone)
	assert.NotNil(t, responseBody.Data.ID)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
}

func TestCreateContactFailed(t *testing.T) {
	TestLogin(t)

	user := new(entity.User)
	err := db.Where("id = ?", "khannedy").First(user).Error
	assert.Nil(t, err)

	requestBody := model.CreateContactRequest{
		FirstName: "",
		LastName:  "",
		Email:     "",
		Phone:     "",
	}
	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	// Buat server Chi untuk testing
	server := httptest.NewServer(app)
	defer server.Close()

	// Kirim request ke server
	req, err := http.NewRequest(http.MethodPost, server.URL+"/api/contacts", strings.NewReader(string(bodyJson)))
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", user.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.ContactResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.NotNil(t, responseBody.Errors)
}

func TestGetContact(t *testing.T) {
	TestCreateContact(t)

	user := new(entity.User)
	err := db.Where("id = ?", "khannedy").First(user).Error
	assert.Nil(t, err)

	contact := new(entity.Contact)
	err = db.Where("user_id = ?", user.ID).First(contact).Error
	assert.Nil(t, err)

	// Buat server Chi untuk testing
	server := httptest.NewServer(app)
	defer server.Close()

	// Kirim request ke server
	req, err := http.NewRequest(http.MethodGet, server.URL+"/api/contacts/"+contact.ID, nil)
	assert.Nil(t, err)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", user.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.ContactResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, contact.ID, responseBody.Data.ID)
	assert.Equal(t, contact.FirstName, responseBody.Data.FirstName)
	assert.Equal(t, contact.LastName, responseBody.Data.LastName)
	assert.Equal(t, contact.Email, responseBody.Data.Email)
	assert.Equal(t, contact.Phone, responseBody.Data.Phone)
	assert.Equal(t, contact.CreatedAt, responseBody.Data.CreatedAt)
	assert.Equal(t, contact.UpdatedAt, responseBody.Data.UpdatedAt)
}

func TestGetContactFailed(t *testing.T) {
	TestCreateContact(t)

	user := new(entity.User)
	err := db.Where("id = ?", "khannedy").First(user).Error
	assert.Nil(t, err)

	// Buat server Chi untuk testing
	server := httptest.NewServer(app)
	defer server.Close()

	// Kirim request ke server
	req, err := http.NewRequest(http.MethodGet, server.URL+"/api/contacts/"+uuid.NewString(), nil)
	assert.Nil(t, err)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", user.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.ContactResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestUpdateContact(t *testing.T) {
	TestCreateContact(t)

	user := new(entity.User)
	err := db.Where("id = ?", "khannedy").First(user).Error
	assert.Nil(t, err)

	contact := new(entity.Contact)
	err = db.Where("user_id = ?", user.ID).First(contact).Error
	assert.Nil(t, err)

	requestBody := model.UpdateContactRequest{
		FirstName: "Eko",
		LastName:  "Budiman",
		Email:     "budiman@example.com",
		Phone:     "089898989",
	}
	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	// Buat server Chi untuk testing
	server := httptest.NewServer(app)
	defer server.Close()

	// Kirim request ke server
	req, err := http.NewRequest(http.MethodPut, server.URL+"/api/contacts/"+contact.ID, strings.NewReader(string(bodyJson)))
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", user.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.ContactResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, requestBody.FirstName, responseBody.Data.FirstName)
	assert.Equal(t, requestBody.LastName, responseBody.Data.LastName)
	assert.Equal(t, requestBody.Email, responseBody.Data.Email)
	assert.Equal(t, requestBody.Phone, responseBody.Data.Phone)
	assert.NotNil(t, responseBody.Data.ID)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
}

func TestUpdateContactFailed(t *testing.T) {
	TestCreateContact(t)

	user := new(entity.User)
	err := db.Where("id = ?", "khannedy").First(user).Error
	assert.Nil(t, err)

	contact := new(entity.Contact)
	err = db.Where("user_id = ?", user.ID).First(contact).Error
	assert.Nil(t, err)

	requestBody := model.UpdateContactRequest{
		FirstName: "",
		LastName:  "",
		Email:     "",
		Phone:     "",
	}
	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	// Buat server Chi untuk testing
	server := httptest.NewServer(app)
	defer server.Close()

	// Kirim request ke server
	req, err := http.NewRequest(http.MethodPut, server.URL+"/api/contacts/"+contact.ID, strings.NewReader(string(bodyJson)))
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", user.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.ContactResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestUpdateContactNotFound(t *testing.T) {
	TestCreateContact(t)

	user := new(entity.User)
	err := db.Where("id = ?", "khannedy").First(user).Error
	assert.Nil(t, err)

	requestBody := model.UpdateContactRequest{
		FirstName: "",
		LastName:  "",
		Email:     "",
		Phone:     "",
	}
	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	// Buat server Chi untuk testing
	server := httptest.NewServer(app)
	defer server.Close()

	// Kirim request ke server
	req, err := http.NewRequest(http.MethodPut, server.URL+"/api/contacts/"+uuid.NewString(), strings.NewReader(string(bodyJson)))
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", user.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.ContactResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestDeleteContact(t *testing.T) {
	TestCreateContact(t)

	user := new(entity.User)
	err := db.Where("id = ?", "khannedy").First(user).Error
	assert.Nil(t, err)

	contact := new(entity.Contact)
	err = db.Where("user_id = ?", user.ID).First(contact).Error
	assert.Nil(t, err)

	// Buat server Chi untuk testing
	server := httptest.NewServer(app)
	defer server.Close()

	// Kirim request ke server
	req, err := http.NewRequest(http.MethodDelete, server.URL+"/api/contacts/"+contact.ID, nil)
	assert.Nil(t, err)
	req.Header.Set("Accept", "application/json")
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
	assert.Equal(t, true, responseBody.Data)
}

func TestDeleteContactFailed(t *testing.T) {
	TestCreateContact(t)

	user := new(entity.User)
	err := db.Where("id = ?", "khannedy").First(user).Error
	assert.Nil(t, err)

	// Buat server Chi untuk testing
	server := httptest.NewServer(app)
	defer server.Close()

	// Kirim request ke server
	req, err := http.NewRequest(http.MethodDelete, server.URL+"/api/contacts/"+uuid.NewString(), nil)
	assert.Nil(t, err)
	req.Header.Set("Accept", "application/json")
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

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestSearchContact(t *testing.T) {
	TestLogin(t)

	user := new(entity.User)
	err := db.Where("id = ?", "khannedy").First(user).Error
	assert.Nil(t, err)

	CreateContacts(user, 20)

	// Buat server Chi untuk testing
	server := httptest.NewServer(app)
	defer server.Close()

	// Kirim request ke server
	req, err := http.NewRequest(http.MethodGet, server.URL+"/api/contacts", nil)
	assert.Nil(t, err)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", user.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[[]model.ContactResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 10, len(responseBody.Data))
	assert.Equal(t, int64(20), responseBody.Paging.TotalItem)
	assert.Equal(t, int64(2), responseBody.Paging.TotalPage)
	assert.Equal(t, 1, responseBody.Paging.Page)
	assert.Equal(t, 10, responseBody.Paging.Size)
}

func TestSearchContactWithPagination(t *testing.T) {
	TestLogin(t)

	user := new(entity.User)
	err := db.Where("id = ?", "khannedy").First(user).Error
	assert.Nil(t, err)

	CreateContacts(user, 20)

	// Buat server Chi untuk testing
	server := httptest.NewServer(app)
	defer server.Close()

	// Kirim request ke server
	req, err := http.NewRequest(http.MethodGet, server.URL+"/api/contacts?page=2&size=5", nil)
	assert.Nil(t, err)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", user.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[[]model.ContactResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 5, len(responseBody.Data))
	assert.Equal(t, int64(20), responseBody.Paging.TotalItem)
	assert.Equal(t, int64(4), responseBody.Paging.TotalPage)
	assert.Equal(t, 2, responseBody.Paging.Page)
	assert.Equal(t, 5, responseBody.Paging.Size)
}

func TestSearchContactWithFilter(t *testing.T) {
	TestLogin(t)

	user := new(entity.User)
	err := db.Where("id = ?", "khannedy").First(user).Error
	assert.Nil(t, err)

	CreateContacts(user, 20)

	// Buat server Chi untuk testing
	server := httptest.NewServer(app)
	defer server.Close()

	// Kirim request ke server
	req, err := http.NewRequest(http.MethodGet, server.URL+"/api/contacts?name=contact&phone=08000000&email=example.com", nil)
	assert.Nil(t, err)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", user.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[[]model.ContactResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 10, len(responseBody.Data))
	assert.Equal(t, int64(20), responseBody.Paging.TotalItem)
	assert.Equal(t, int64(2), responseBody.Paging.TotalPage)
	assert.Equal(t, 1, responseBody.Paging.Page)
	assert.Equal(t, 10, responseBody.Paging.Size)
}