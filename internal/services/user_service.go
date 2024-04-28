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

/*
   UserService
       Constructor() (Wajib)
       RegisterRoutes() (Wajib)

       handleUserRegister()
       handleUserLogin()
*/
type UserService struct {
	Store datastore.Store
}

// konstruktor untuk user service
func NewUserService(s datastore.Store) *UserService {

	return &UserService{
		Store: s,
	}
}

func (s *UserService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/user/register", helper.CreateHandlerFunc(s.handleUserRegister)).Methods(http.MethodPost)
	r.HandleFunc("/user/login", helper.CreateHandlerFunc(s.handleUserLogin)).Methods(http.MethodPost)

	r.HandleFunc("/user/{id}", helper.CreateHandlerFunc(auth.JWTMiddleware(s.handleUserUpdate))).Methods(http.MethodPatch)
	r.HandleFunc("/user/{id}", helper.CreateHandlerFunc(auth.JWTMiddleware(s.handleUserDelete))).Methods(http.MethodDelete)
}

func (s *UserService) handleUserUpdate(w http.ResponseWriter, r *http.Request) types.AppError {

	err := usecases.UpdateUser(s.Store, w, r)
	if err.Error != nil {
		return err
	}

	return types.AppError{
		Error:  nil,
		Status: http.StatusOK,
	}
}

func (s *UserService) handleUserDelete(w http.ResponseWriter, r *http.Request) types.AppError {

	err := usecases.DeleteUser(s.Store, w, r)
	if err.Error != nil {
		return err
	}

	return types.AppError{
		Error:  nil,
		Status: http.StatusOK,
	}
}

func (s *UserService) handleUserRegister(w http.ResponseWriter, r *http.Request) types.AppError {

	err := usecases.CreateUser(s.Store, w, r)
	if err.Error != nil {
		return err
	}

	return types.AppError{
		Error:  nil,
		Status: http.StatusCreated,
	}
}

func (s *UserService) handleUserLogin(w http.ResponseWriter, r *http.Request) types.AppError {

	err := usecases.AuthorizeUser(s.Store, w, r)
	if err.Error != nil {

		return err
	}

	return types.AppError{
		Error:  nil,
		Status: http.StatusOK,
	}
}
