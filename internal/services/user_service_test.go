package services

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GetterSethya/golangApiMarketplace/internal/datastore"
	"github.com/GetterSethya/golangApiMarketplace/internal/entities"
	"github.com/GetterSethya/golangApiMarketplace/internal/helper"
	"github.com/gorilla/mux"
)

func TestCreateUser(t *testing.T) {
	inMemoryDb := datastore.MockStore{}
	userService := NewUserService(&inMemoryDb)

	t.Run("Should return an error if name is empty", func(t *testing.T) {
		payload := &entities.User{
			Name: "",
		}

		b, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(b))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/user", helper.CreateHandlerFunc(userService.handleUserRegister))
		router.ServeHTTP(rr, req)

		expectedCode := http.StatusBadRequest

		if rr.Code != expectedCode {

			t.Errorf("Invalid status code, expected: %d, but got: %d", expectedCode, rr.Code)
		}
	})

	t.Run("Should return an error if username is empty", func(t *testing.T) {
		payload := &entities.User{
			Username: "",
		}

		b, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(b))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/user", helper.CreateHandlerFunc(userService.handleUserRegister))
		router.ServeHTTP(rr, req)

		expectedCode := http.StatusBadRequest

		if rr.Code != expectedCode {

			t.Errorf("Invalid status code, expected: %d, but got: %d", expectedCode, rr.Code)
		}
	})

	t.Run("Should return an error if password is empty", func(t *testing.T) {
		payload := &entities.User{
			HashPassword: "",
		}

		b, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(b))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/user", helper.CreateHandlerFunc(userService.handleUserRegister))
		router.ServeHTTP(rr, req)

		expectedCode := http.StatusBadRequest

		if rr.Code != expectedCode {

			t.Errorf("Invalid status code, expected: %d, but got: %d", expectedCode, rr.Code)
		}
	})

	t.Run("Should return an error if username has space character", func(t *testing.T) {
		payload := &entities.User{
			Username: "user name",
		}

		b, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(b))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/user", helper.CreateHandlerFunc(userService.handleUserRegister))
		router.ServeHTTP(rr, req)

		expectedCode := http.StatusBadRequest

		if rr.Code != expectedCode {

			t.Errorf("Invalid status code, expected: %d, but got: %d", expectedCode, rr.Code)
		}
	})

	t.Run("Should create user", func(t *testing.T) {
		payload := &entities.User{
			Username:     "username",
			Name:         "nama user",
			HashPassword: "12345678",
		}

		b, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(b))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/user", helper.CreateHandlerFunc(userService.handleUserRegister))
		router.ServeHTTP(rr, req)

		expectedCode := http.StatusCreated

		if rr.Code != expectedCode {

			t.Errorf("Invalid status code, expected: %d, but got: %d", expectedCode, rr.Code)
		}
	})
}
