package test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/iyasz/golang-clean-architecture/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestCreateAddress(t *testing.T) {
	TestCreateContact(t)

	user := GetFirstUser(t)
	contact := GetFirstContact(t, user)

	requestBody := model.CreateAddressRequest{
		Street:     "Jalan Belum Jadi",
		City:       "Jakarta",
		Province:   "DKI Jakarta",
		PostalCode: "343443",
		Country:    "Indonesia",
	}
	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	server := httptest.NewServer(app)
	defer server.Close()

	req, err := http.NewRequest(http.MethodPost, server.URL+BaseContactsAPIURL+"/"+contact.ID+"/addresses", strings.NewReader(string(bodyJson)))
	assert.Nil(t, err)
	SetupHeader(req)
	req.Header.Set("Authorization", user.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.AddressResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, requestBody.Street, responseBody.Data.Street)
	assert.Equal(t, requestBody.City, responseBody.Data.City)
	assert.Equal(t, requestBody.Province, responseBody.Data.Province)
	assert.Equal(t, requestBody.Country, responseBody.Data.Country)
	assert.Equal(t, requestBody.PostalCode, responseBody.Data.PostalCode)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
	assert.NotNil(t, responseBody.Data.ID)
}

func TestCreateAddressFailed(t *testing.T) {
	TestCreateContact(t)

	user := GetFirstUser(t)
	contact := GetFirstContact(t, user)

	requestBody := model.CreateAddressRequest{
		Street:     "Jalan Belum Jadi",
		City:       "Jakarta",
		Province:   "DKI Jakarta",
		PostalCode: "343443343443343443343443343443343443343443",
		Country:    "Indonesia",
	}
	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	server := httptest.NewServer(app)
	defer server.Close()

	req, err := http.NewRequest(http.MethodPost, server.URL+BaseContactsAPIURL+"/"+contact.ID+"/addresses", strings.NewReader(string(bodyJson)))
	assert.Nil(t, err)
	SetupHeader(req)
	req.Header.Set("Authorization", user.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.AddressResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestListAddresses(t *testing.T) {
	TestCreateContact(t)

	user := GetFirstUser(t)
	contact := GetFirstContact(t, user)

	CreateAddresses(t, contact, 5)

	server := httptest.NewServer(app)
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL+BaseContactsAPIURL+"/"+contact.ID+"/addresses", nil)
	assert.Nil(t, err)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", user.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[[]model.AddressResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 5, len(responseBody.Data))
}

func TestListAddressesFailed(t *testing.T) {
	TestCreateContact(t)

	user := GetFirstUser(t)
	contact := GetFirstContact(t, user)

	CreateAddresses(t, contact, 5)

	server := httptest.NewServer(app)
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL+BaseContactsAPIURL+"/"+"wrong"+"/addresses", nil)
	assert.Nil(t, err)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", user.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[[]model.AddressResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestGetAddress(t *testing.T) {
	TestCreateAddress(t)

	user := GetFirstUser(t)
	contact := GetFirstContact(t, user)
	address := GetFirstAddress(t, contact)

	server := httptest.NewServer(app)
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL+BaseContactsAPIURL+"/"+contact.ID+"/addresses/"+address.ID, nil)
	assert.Nil(t, err)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", user.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.AddressResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, address.ID, responseBody.Data.ID)
	assert.Equal(t, address.Street, responseBody.Data.Street)
	assert.Equal(t, address.City, responseBody.Data.City)
	assert.Equal(t, address.Province, responseBody.Data.Province)
	assert.Equal(t, address.Country, responseBody.Data.Country)
	assert.Equal(t, address.PostalCode, responseBody.Data.PostalCode)
	assert.Equal(t, address.CreatedAt, responseBody.Data.CreatedAt)
	assert.Equal(t, address.UpdatedAt, responseBody.Data.UpdatedAt)
}

func TestGetAddressFailed(t *testing.T) {
	TestCreateAddress(t)

	user := GetFirstUser(t)
	contact := GetFirstContact(t, user)

	server := httptest.NewServer(app)
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL+BaseContactsAPIURL+"/"+contact.ID+"/addresses/"+"wrong", nil)
	assert.Nil(t, err)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", user.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.AddressResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestUpdateAddress(t *testing.T) {
	TestCreateAddress(t)

	user := GetFirstUser(t)
	contact := GetFirstContact(t, user)
	address := GetFirstAddress(t, contact)

	requestBody := model.CreateAddressRequest{
		Street:     "Jalan Lagi Dijieun",
		City:       "Bandung",
		Province:   "Jawa Barat",
		PostalCode: "343443",
		Country:    "Indonesia",
	}
	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	server := httptest.NewServer(app)
	defer server.Close()

	req, err := http.NewRequest(http.MethodPut, server.URL+BaseContactsAPIURL+"/"+contact.ID+"/addresses/"+address.ID, strings.NewReader(string(bodyJson)))
	assert.Nil(t, err)
	SetupHeader(req)
	req.Header.Set("Authorization", user.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.AddressResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, requestBody.Street, responseBody.Data.Street)
	assert.Equal(t, requestBody.City, responseBody.Data.City)
	assert.Equal(t, requestBody.Province, responseBody.Data.Province)
	assert.Equal(t, requestBody.Country, responseBody.Data.Country)
	assert.Equal(t, requestBody.PostalCode, responseBody.Data.PostalCode)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
	assert.NotNil(t, responseBody.Data.ID)
}

func TestUpdateAddressFailed(t *testing.T) {
	TestCreateAddress(t)

	user := GetFirstUser(t)
	contact := GetFirstContact(t, user)
	address := GetFirstAddress(t, contact)

	requestBody := model.UpdateAddressRequest{
		Street:     "Jalan Lagi Dijieun",
		City:       "Bandung",
		Province:   "Jawa Barat",
		PostalCode: "343443343443343443343443343443343443343443",
		Country:    "Indonesia",
	}
	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	server := httptest.NewServer(app)
	defer server.Close()

	req, err := http.NewRequest(http.MethodPut, server.URL+BaseContactsAPIURL+"/"+contact.ID+"/addresses/"+address.ID, strings.NewReader(string(bodyJson)))
	assert.Nil(t, err)
	SetupHeader(req)
	req.Header.Set("Authorization", user.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.AddressResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestDeleteAddress(t *testing.T) {
	TestCreateAddress(t)

	user := GetFirstUser(t)
	contact := GetFirstContact(t, user)
	address := GetFirstAddress(t, contact)

	server := httptest.NewServer(app)
	defer server.Close()

	req, err := http.NewRequest(http.MethodDelete, server.URL+BaseContactsAPIURL+"/"+contact.ID+"/addresses/"+address.ID, nil)
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

func TestDeleteAddressFailed(t *testing.T) {
	TestCreateAddress(t)

	user := GetFirstUser(t)
	contact := GetFirstContact(t, user)

	server := httptest.NewServer(app)
	defer server.Close()

	// Kirim request to server
	req, err := http.NewRequest(http.MethodDelete, server.URL+BaseContactsAPIURL+"/"+contact.ID+"/addresses/"+"wrong", nil)
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