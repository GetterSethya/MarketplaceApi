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

type Transactionservice struct {
	Store datastore.Store
}

func NewTransactionService(s datastore.Store) *Transactionservice {

	return &Transactionservice{
		Store: s,
	}
}

func (s *Transactionservice) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/transaction/{id}", helper.CreateHandlerFunc(auth.JWTMiddleware(s.GetTransaction))).Methods(http.MethodGet)
	r.HandleFunc("/transaction/{id}", helper.CreateHandlerFunc(auth.JWTMiddleware(s.UpdateStatusTransaction))).Methods(http.MethodPatch)
	r.HandleFunc("/transaction", helper.CreateHandlerFunc(auth.JWTMiddleware(s.ListTransaction))).Methods(http.MethodGet)
	r.HandleFunc("/transaction", helper.CreateHandlerFunc(auth.JWTMiddleware(s.CreateTransaction))).Methods(http.MethodPost)
}

func (s *Transactionservice) UpdateStatusTransaction(w http.ResponseWriter, r *http.Request) types.AppError {

	if err := usecases.UpdateStatusTransaction(s.Store, w, r); err.Error != nil {
		return err
	}

	return types.AppError{
		Error:  nil,
		Status: http.StatusOK,
	}
}

func (s *Transactionservice) CreateTransaction(w http.ResponseWriter, r *http.Request) types.AppError {

	if err := usecases.CreateTransaction(s.Store, w, r); err.Error != nil {
		return err
	}

	return types.AppError{
		Error:  nil,
		Status: http.StatusCreated,
	}
}

func (s *Transactionservice) ListTransaction(w http.ResponseWriter, r *http.Request) types.AppError {

	if err := usecases.ListTransaction(s.Store, w, r); err.Error != nil {
		return err
	}

	return types.AppError{
		Error:  nil,
		Status: http.StatusOK,
	}
}

func (s *Transactionservice) GetTransaction(w http.ResponseWriter, r *http.Request) types.AppError {

	if err := usecases.GetTransaction(s.Store, w, r); err.Error != nil {
		return err
	}

	return types.AppError{
		Error:  nil,
		Status: http.StatusOK,
	}
}
