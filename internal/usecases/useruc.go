package usecases

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"

	"github.com/GetterSethya/golangApiMarketplace/config"
	"github.com/GetterSethya/golangApiMarketplace/internal/auth"
	"github.com/GetterSethya/golangApiMarketplace/internal/datastore"
	"github.com/GetterSethya/golangApiMarketplace/internal/entities"
	"github.com/GetterSethya/golangApiMarketplace/internal/helper"
	"github.com/GetterSethya/golangApiMarketplace/internal/types"
	"github.com/GetterSethya/golangApiMarketplace/internal/validator"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type UserUseCase interface {
	CreateUser(s datastore.Store, w http.ResponseWriter, r *http.Request) types.AppError
	GetUserById(s datastore.Store, w http.ResponseWriter, r *http.Request) (*entities.User, types.AppError)
	GetUserByUsername(s datastore.Store, w http.ResponseWriter, r *http.Request) (*entities.User, types.AppError)
	UpdateUser(s datastore.Store, w http.ResponseWriter, r *http.Request) (*entities.User, types.AppError)
	DeleteUser(s datastore.Store, w http.ResponseWriter, r *http.Request) types.AppError
	AuthorizeUser(s datastore.Store, w http.ResponseWriter, r *http.Request) types.AppError
}

func CreateUser(s datastore.Store, w http.ResponseWriter, r *http.Request) types.AppError {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error when reading body")

		return types.AppError{
			Error:  fmt.Errorf("Invalid username/name/password"),
			Status: http.StatusBadRequest,
		}
	}

	defer r.Body.Close()

	var user *entities.User

	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Println("Error when unmarshaling body")

		return types.AppError{
			Error: fmt.Errorf("Invalid username/name/password"),
		}
	}

	if err := validator.ValidateRegisterPayload(user); err != nil {

		log.Println("Error when validating user from input body:", err)

		return types.AppError{
			Error:  fmt.Errorf(err.Error()),
			Status: http.StatusBadRequest,
		}
	}

	id := uuid.NewString()

	if err := s.CreateUser(id, user); err != nil {
		log.Println("Error when creating user in userUseCase:", err)

		if reflect.TypeOf(err).String() == "*pq.Error" {

			return types.AppError{
				Error:  fmt.Errorf("Failed when registering user, username already taken"),
				Status: http.StatusNotAcceptable,
			}
		}

		return types.AppError{
			Error:  fmt.Errorf("Failed when registering user, please try again"),
			Status: http.StatusInternalServerError,
		}
	}

	secret := config.LoadConfig().App.JWTSecret
	accessToken, err := auth.CreateJWT(id, secret)
	if err != nil {

		log.Println("Error when creating JWT")

		return types.AppError{
			Error:  fmt.Errorf("Failed when registering user"),
			Status: http.StatusInternalServerError,
		}
	}

	resp := types.ServerResponse{
		Message: "User registered successfully",
		Data: map[string]interface{}{
			"username":    user.Username,
			"name":        user.Name,
			"accessToken": accessToken,
		},
	}

	helper.WriteJson(w, http.StatusCreated, resp)

	return types.AppError{
		Error:  nil,
		Status: 201,
	}
}

func AuthorizeUser(s datastore.Store, w http.ResponseWriter, r *http.Request) types.AppError {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error when reading body")

		return types.AppError{
			Error:  fmt.Errorf("Invalid username/password"),
			Status: http.StatusBadRequest,
		}
	}

	defer r.Body.Close()
	if string(body) == "" {
		log.Println("error body is nil")

		return types.AppError{
			Error:  fmt.Errorf("Invalid username/name/password"),
			Status: http.StatusBadRequest,
		}
	}

	var user *entities.User

	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Println("Error when unmarshaling body")

		return types.AppError{
			Error: fmt.Errorf("Invalid username/password"),
		}
	}

	if err := validator.ValidateLoginPayload(user); err != nil {
		log.Println("Error when validating user from input body:", err)

		return types.AppError{
			Error:  fmt.Errorf(err.Error()),
			Status: http.StatusBadRequest,
		}
	}

	user, err = s.GetUserByUsername(user.Username)
	if err != nil {
		log.Println("Error when in GetUserByUsername in useruc.go:", err)

		return types.AppError{
			Error:  fmt.Errorf(err.Error()),
			Status: http.StatusNotFound,
		}
	}

	jwtToken, err := auth.CreateJWT(user.ID, config.LoadConfig().App.JWTSecret)
	if err != nil {
		log.Println("Error when creating JWT in useruc.go:", err)

		return types.AppError{
			Error:  fmt.Errorf("Something went wrong, please try again"),
			Status: http.StatusInternalServerError,
		}
	}

	resp := types.ServerResponse{
		Message: "Login succesfull",
		Data: map[string]interface{}{
			"username":    user.Username,
			"name":        user.Name,
			"accessToken": jwtToken,
		},
	}

	helper.WriteJson(w, http.StatusOK, resp)

	return types.AppError{
		Error:  nil,
		Status: 200,
	}
}

func GetUserById(s datastore.Store, w http.ResponseWriter, r *http.Request) (*entities.User, types.AppError) {
	//TODO
	// fungsi s.GetUserById(id) udah ada, tinggal ambil id dari path -> validasi -> GetUserById() -> return json

	return &entities.User{}, types.AppError{
		Error:  nil,
		Status: 200,
	}
}

func GetUserByUsername(s datastore.Store, w http.ResponseWriter, r *http.Request) (*entities.User, types.AppError) {
	//TODO
	// fungsi s.GetUserByUsername(username) udah ada, tinggal ambil username dari path -> validasi -> GetUserByUsername() -> return json

	return &entities.User{}, types.AppError{
		Error:  nil,
		Status: 200,
	}
}

func UpdateUser(s datastore.Store, w http.ResponseWriter, r *http.Request) types.AppError {
	vars := mux.Vars(r)
	userIdUrlPath := vars["id"]
	userIdJWT := auth.GetUserIdFromJWT(r)

	if userIdJWT != userIdUrlPath {

		return types.AppError{
			Error:  fmt.Errorf("Forbidden"),
			Status: 403,
		}
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error when reading body")

		return types.AppError{
			Error:  fmt.Errorf("Invalid username/password"),
			Status: http.StatusBadRequest,
		}
	}

	defer r.Body.Close()

	var user *entities.User

	err = json.Unmarshal(body, &user)

	if err != nil {
		log.Println("Error when unmarshaling body", err)

		return types.AppError{
			Error:  fmt.Errorf("Invalid payload"),
			Status: http.StatusBadRequest,
		}
	}

	if err := validator.ValidateUpdatePayload(user); err != nil {

		log.Println("Error when validating update payload", err)

		return types.AppError{
			Error:  err,
			Status: http.StatusBadRequest,
		}
	}

	if err := s.UpdateUser(userIdJWT, user.Name, user.Username); err != nil {

		log.Println("Error when updating user in useruc.go", err)

		return types.AppError{
			Error:  fmt.Errorf("Something went wrong. Please try again"),
			Status: http.StatusInternalServerError,
		}
	}

	user, err = s.GetUserById(userIdJWT)
	if err != nil {

		log.Println("error when getting user by in ind useruc.go", err)

		return types.AppError{
			Error:  fmt.Errorf("Something went wrong. Please try again"),
			Status: http.StatusInternalServerError,
		}
	}

	resp := types.ServerResponse{
		Message: "User edited successfully",
		Data: map[string]interface{}{
			"user": user,
		},
	}

	helper.WriteJson(w, http.StatusOK, resp)

	return types.AppError{
		Error:  nil,
		Status: http.StatusOK,
	}
}

func DeleteUser(s datastore.Store, w http.ResponseWriter, r *http.Request) types.AppError {

	vars := mux.Vars(r)
	userIdUrlPath := vars["id"]
	userIdJWT := auth.GetUserIdFromJWT(r)

	if userIdJWT != userIdUrlPath {

		return types.AppError{
			Error:  fmt.Errorf("Forbidden"),
			Status: 403,
		}
	}

	err := s.DeleteUser(userIdJWT)
	if err != nil {
		return types.AppError{
			Error:  fmt.Errorf("Failed when deleting user/user didnot exists, please try again"),
			Status: http.StatusInternalServerError,
		}
	}

	resp := types.ServerResponse{
		Message: "User deleted successfully",
		Data:    nil,
	}

	helper.WriteJson(w, http.StatusOK, resp)

	return types.AppError{
		Error:  nil,
		Status: 200,
	}
}
