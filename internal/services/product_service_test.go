package services

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/GetterSethya/golangApiMarketplace/internal/auth"
	"github.com/GetterSethya/golangApiMarketplace/internal/datastore"
	"github.com/GetterSethya/golangApiMarketplace/internal/entities"
	"github.com/GetterSethya/golangApiMarketplace/internal/helper"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func TestCreateProduct(t *testing.T) {
	inMemoryDb := datastore.MockStore{}
	productService := NewProductService(&inMemoryDb)

	t.Run("Should create product", func(t *testing.T) {
		payload := &entities.Product{
			Name:           "nama produk",
			Price:          15000,
			ImageUrl:       "asoidsdas",
			Stock:          10,
			Condition:      "new",
			Tags:           []string{"#murah", "#mantap"},
			IsPurchaseable: true,
			SellerId:       "1923810293812903",
		}

		b, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/product", bytes.NewBuffer(b))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/product", helper.CreateHandlerFunc(productService.handleCreateProduct))
		router.ServeHTTP(rr, req)

		expectedCode := http.StatusCreated

		if rr.Code != expectedCode {

			t.Errorf("Invalid status code, expected: %d, but got: %d", expectedCode, rr.Code)
		}

	})
}

func TestGetProduct(t *testing.T) {

	inMemoryDb := datastore.MockStore{}
	productService := NewProductService(&inMemoryDb)

	t.Run("Should get product", func(t *testing.T) {

		req, err := http.NewRequest(http.MethodGet, "/product/b78cd7e2-765e-4344-aa83-9b61aaa3dec4", bytes.NewBuffer([]byte{}))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/product/{id}", helper.CreateHandlerFunc(productService.handleGetProduct)).Methods(http.MethodGet)
		router.ServeHTTP(rr, req)

		expectedCode := http.StatusOK

		if rr.Code != expectedCode {

			t.Errorf("Invalid status code, expected: %d, but got: %d", expectedCode, rr.Code)
		}

	})

	t.Run("Should 404 when id is not uuid", func(t *testing.T) {

		req, err := http.NewRequest(http.MethodGet, "/product/1234785981bnsakdjnf", bytes.NewBuffer([]byte{}))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/product/{id}", helper.CreateHandlerFunc(productService.handleGetProduct)).Methods(http.MethodGet)
		router.ServeHTTP(rr, req)

		expectedCode := http.StatusNotFound

		if rr.Code != expectedCode {

			t.Errorf("Invalid status code, expected: %d, but got: %d", expectedCode, rr.Code)
		}
	})
}

func TestEditProduct(t *testing.T) {

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	inMemoryDb := datastore.MockStore{}
	productService := NewProductService(&inMemoryDb)
	userId := "75ea96d2-8077-48aa-aad6-a02fbd282f3c"

	t.Run("Should edit product", func(t *testing.T) {
		payload := &entities.Product{
			Name:           "nama produk",
			Price:          15000,
			ImageUrl:       "asoidsdas",
			Stock:          10,
			Condition:      "new",
			Tags:           []string{"#murah", "#mantap"},
			IsPurchaseable: true,
			SellerId:       userId,
		}

		token, err := auth.CreateJWT(userId, "qnqwienidbfsldjlsdf")
		if err != nil {
			t.Error(err)
		}

		b, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/product/b78cd7e2-765e-4344-aa83-9b61aaa3dec4", bytes.NewBuffer(b))
		req.Header.Add("authorization", token)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/product/{id}", helper.CreateHandlerFunc(productService.handleUpdateProduct)).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		expectedCode := http.StatusOK

		if rr.Code != expectedCode {

			t.Errorf("Invalid status code, expected: %d, but got: %d", expectedCode, rr.Code)
		}

	})
}
