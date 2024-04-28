package services

import (
	"net/http"

	"github.com/GetterSethya/golangApiMarketplace/internal/auth"
	"github.com/GetterSethya/golangApiMarketplace/internal/datastore"
	"github.com/GetterSethya/golangApiMarketplace/internal/helper"
	"github.com/GetterSethya/golangApiMarketplace/internal/types"
	"github.com/GetterSethya/golangApiMarketplace/internal/usecases"
	"github.com/gorilla/mux"
)

type ProductService struct {
	Store datastore.Store
}

func NewProductService(s datastore.Store) *ProductService {

	return &ProductService{
		Store: s,
	}
}

func (s *ProductService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/product", helper.CreateHandlerFunc(auth.JWTMiddleware(s.handleCreateProduct))).Methods(http.MethodPost)
	r.HandleFunc("/product/{id}", helper.CreateHandlerFunc(auth.JWTMiddleware(s.handleUpdateProduct))).Methods(http.MethodPatch)
	r.HandleFunc("/product/{id}", helper.CreateHandlerFunc(s.handleGetProduct)).Methods(http.MethodGet)
	r.HandleFunc("/product", helper.CreateHandlerFunc(s.handleListProduct)).Methods(http.MethodGet)
	r.HandleFunc("/product/{id}", helper.CreateHandlerFunc(auth.JWTMiddleware(s.handleDeleteProduct))).Methods(http.MethodDelete)
	r.HandleFunc("/product/{id}/stock", helper.CreateHandlerFunc(auth.JWTMiddleware(s.handleUpdateStock))).Methods(http.MethodPost)
}

func (s *ProductService) handleUpdateStock(w http.ResponseWriter, r *http.Request) types.AppError {

	if err := usecases.UpdateStock(s.Store, w, r); err.Error != nil {
		return err
	}

	return types.AppError{
		Error:  nil,
		Status: http.StatusOK,
	}
}

func (s *ProductService) handleDeleteProduct(w http.ResponseWriter, r *http.Request) types.AppError {

	if err := usecases.DeleteProduct(s.Store, w, r); err.Error != nil {

		return err
	}

	return types.AppError{
		Error:  nil,
		Status: http.StatusOK,
	}
}

func (s *ProductService) handleListProduct(w http.ResponseWriter, r *http.Request) types.AppError {

	if err := usecases.ListProduct(s.Store, w, r); err.Error != nil {

		return err
	}

	return types.AppError{
		Error:  nil,
		Status: http.StatusOK,
	}
}

func (s *ProductService) handleGetProduct(w http.ResponseWriter, r *http.Request) types.AppError {

	// ini nampilin product detail, GET /v1/product/{id}
	if err := usecases.GetProductById(s.Store, w, r); err.Error != nil {
		return err
	}

	return types.AppError{
		Error:  nil,
		Status: http.StatusOK,
	}
}

func (s *ProductService) handleUpdateProduct(w http.ResponseWriter, r *http.Request) types.AppError {

	if err := usecases.UpdateProduct(s.Store, w, r); err.Error != nil {

		return err
	}

	return types.AppError{
		Error:  nil,
		Status: http.StatusOK,
	}
}

func (s *ProductService) handleCreateProduct(w http.ResponseWriter, r *http.Request) types.AppError {

	if err := usecases.CreateProduct(s.Store, w, r); err.Error != nil {

		return err
	}

	return types.AppError{
		Error:  nil,
		Status: http.StatusCreated,
	}
}
