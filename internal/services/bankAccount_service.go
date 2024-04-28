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

type BankAccountService struct {
	Store datastore.Store
}


func NewBankAccountService(s datastore.Store) *BankAccountService {

	return &BankAccountService{
		Store: s,
	}
}

func (s *BankAccountService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/bank/account/user/{id}", helper.CreateHandlerFunc(s.handleListBankAccount)).Methods(http.MethodGet)
	r.HandleFunc("/bank/account", helper.CreateHandlerFunc(auth.JWTMiddleware(s.handleCreateBankAccount))).Methods(http.MethodPost)
	r.HandleFunc("/bank/account/{id}", helper.CreateHandlerFunc(auth.JWTMiddleware(s.handleUpdateBankAccount))).Methods(http.MethodPatch)
	r.HandleFunc("/bank/account/{id}", helper.CreateHandlerFunc(auth.JWTMiddleware(s.handleDeleteBankAccount))).Methods(http.MethodDelete)
}

func (s *BankAccountService) handleUpdateBankAccount(w http.ResponseWriter, r *http.Request) types.AppError {

	if err := usecases.UpdateBankAccount(s.Store, w, r); err.Error != nil {
		return err
	}

	return types.AppError{
		Error:  nil,
		Status: http.StatusOK,
	}
}

func (s *BankAccountService) handleDeleteBankAccount(w http.ResponseWriter, r *http.Request) types.AppError {

	if err := usecases.DeleteBankAccount(s.Store, w, r); err.Error != nil {
		return err
	}

	return types.AppError{
		Error:  nil,
		Status: http.StatusOK,
	}
}

func (s *BankAccountService) handleListBankAccount(w http.ResponseWriter, r *http.Request) types.AppError {

	if err := usecases.ListBankAccount(s.Store, w, r); err.Error != nil {
		return err
	}

	return types.AppError{
		Error:  nil,
		Status: http.StatusOK,
	}
}

func (s *BankAccountService) handleCreateBankAccount(w http.ResponseWriter, r *http.Request) types.AppError {

	if err := usecases.CreateBankAccount(s.Store, w, r); err.Error != nil {
		return err
	}

	return types.AppError{
		Error:  nil,
		Status: http.StatusCreated,
	}
}
