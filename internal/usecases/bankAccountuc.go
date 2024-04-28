package usecases

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/GetterSethya/golangApiMarketplace/internal/auth"
	"github.com/GetterSethya/golangApiMarketplace/internal/datastore"
	"github.com/GetterSethya/golangApiMarketplace/internal/entities"
	"github.com/GetterSethya/golangApiMarketplace/internal/helper"
	"github.com/GetterSethya/golangApiMarketplace/internal/types"
	"github.com/GetterSethya/golangApiMarketplace/internal/validator"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type BankAccountUseCase interface {
	CreateBankAccount(s datastore.Store, w http.ResponseWriter, r *http.Request) types.AppError
	ListBankAccount(s datastore.Store, w http.ResponseWriter, r *http.Request) types.AppError
	UpdateBankAccount(s datastore.Store, w http.ResponseWriter, r *http.Request) types.AppError
	DeleteBankAccount(s datastore.Store, w http.ResponseWriter, r *http.Request) types.AppError
}

func UpdateBankAccount(s datastore.Store, w http.ResponseWriter, r *http.Request) types.AppError {

	vars := mux.Vars(r)
	bankAccId := vars["id"]
	sellerId := auth.GetUserIdFromJWT(r)

	bankAcc, err := s.GetBankAccount(bankAccId)
	if err != nil {

		log.Println("error when getting bank account in UpdateBankAccount", err)

		return types.AppError{
			Error:  fmt.Errorf("Failed when deleting bank account, please try again."),
			Status: http.StatusInternalServerError,
		}
	}

	if bankAcc.SellerId != sellerId {
		return types.AppError{
			Error:  fmt.Errorf("Forbidden"),
			Status: http.StatusForbidden,
		}
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {

		log.Println("error when reading body from request in UpdateBankAccount in bankAccountuc.go")

		return types.AppError{
			Error:  fmt.Errorf("Invalid/missing field"),
			Status: http.StatusBadRequest,
		}
	}

	defer r.Body.Close()

	var bankAccount *entities.BankAccount

	err = json.Unmarshal(body, &bankAccount)
	if err != nil {

		log.Println("Error when Unmarshal body in update bank account usecases ", err)

		return types.AppError{
			Error:  fmt.Errorf("Invalid/missing field"),
			Status: http.StatusBadRequest,
		}
	}

	if err := validator.ValidateUpdateBankAccountPayload(bankAccount); err != nil {

		log.Println("error when validating update bank account payload", err)

		return types.AppError{
			Error:  err,
			Status: http.StatusBadRequest,
		}
	}

	if err := s.UpdateBankAccount(bankAccId, bankAccount); err != nil {

		log.Println("error when updating bank account in bankaccountuc.go:", err)

		return types.AppError{
			Error:  fmt.Errorf("Failed updating bank account informations"),
			Status: http.StatusInternalServerError,
		}
	}

	respBankAccount, err := s.GetBankAccount(bankAccId)
	if err != nil {
		log.Println("error when getting bankAccount in bankAccountuc.go")

		return types.AppError{
			Error:  fmt.Errorf("Something went wrong when updating bank account"),
			Status: http.StatusInternalServerError,
		}
	}

	resp := types.ServerResponse{
		Message: "Bank account updated successfully",
		Data: map[string]interface{}{
			"bankAccount": respBankAccount,
		},
	}

	helper.WriteJson(w, http.StatusOK, resp)

	return types.AppError{
		Error:  nil,
		Status: http.StatusOK,
	}
}

func DeleteBankAccount(s datastore.Store, w http.ResponseWriter, r *http.Request) types.AppError {

	vars := mux.Vars(r)
	bankAccId := vars["id"]
	sellerId := auth.GetUserIdFromJWT(r)

	bankAcc, err := s.GetBankAccount(bankAccId)
	if err != nil {

		log.Println("error when getting bank account in DeleteBankAccount", err)

		return types.AppError{
			Error:  fmt.Errorf("Failed when deleting bank account, please try again."),
			Status: http.StatusInternalServerError,
		}
	}

	if bankAcc.SellerId != sellerId {
		return types.AppError{
			Error:  fmt.Errorf("Forbidden"),
			Status: http.StatusForbidden,
		}
	}

	err = s.DeleteBankAccount(bankAccId)

	if err != nil {

		log.Println("error when deleting bank account", err)

		return types.AppError{
			Error:  fmt.Errorf("Failed when deleting bank account, please try again."),
			Status: http.StatusInternalServerError,
		}
	}

	resp := types.ServerResponse{
		Message: "Ok",
		Data:    nil,
	}

	helper.WriteJson(w, http.StatusOK, resp)

	return types.AppError{
		Error:  nil,
		Status: http.StatusOK,
	}
}

func ListBankAccount(s datastore.Store, w http.ResponseWriter, r *http.Request) types.AppError {

	vars := mux.Vars(r)
	userIdUrlPath := vars["id"]

	listBankAcc, err := s.ListBankAccount(userIdUrlPath)

	if err != nil {

		log.Println("error when getting list bank account", err)

		return types.AppError{
			Error:  fmt.Errorf("Failed when getting list bank account, please try again."),
			Status: http.StatusInternalServerError,
		}
	}

	resp := types.ServerResponse{
		Message: "Ok",
		Data: map[string]interface{}{
			"bankAccounts": listBankAcc,
		},
	}

	helper.WriteJson(w, http.StatusOK, resp)

	return types.AppError{
		Error:  nil,
		Status: http.StatusOK,
	}
}

func CreateBankAccount(s datastore.Store, w http.ResponseWriter, r *http.Request) types.AppError {

	sellerId := auth.GetUserIdFromJWT(r)

	body, err := io.ReadAll(r.Body)
	if err != nil {

		log.Println("Error when reading body in CreateBankAccount usecase")

		return types.AppError{
			Error:  fmt.Errorf("Invalid/missing field"),
			Status: http.StatusBadRequest,
		}
	}

	defer r.Body.Close()

	var bankAccount *entities.BankAccount

	err = json.Unmarshal(body, &bankAccount)
	if err != nil {

		log.Println("error when Unmarshal body in create product usecase", err)

		return types.AppError{
			Error:  fmt.Errorf("Invalid/missing field"),
			Status: http.StatusBadRequest,
		}
	}

	if err := validator.ValidateCreateBankAccountPayload(bankAccount); err != nil {

		log.Println("error when validating create product payload")

		return types.AppError{
			Error:  err,
			Status: http.StatusBadRequest,
		}
	}

	id := uuid.NewString()

	if err := s.CreateBankAccount(id, sellerId, bankAccount); err != nil {

		log.Println("error when creating product", err)

		return types.AppError{
			Error:  fmt.Errorf("Failed when creating product, please try again."),
			Status: http.StatusInternalServerError,
		}
	}

	newBankAccount, err := s.GetBankAccount(id)

	if err != nil {

		log.Println("error when getting bank account", err)

		return types.AppError{
			Error:  fmt.Errorf("Failed when creating bank account, please try again."),
			Status: http.StatusInternalServerError,
		}
	}

	resp := types.ServerResponse{
		Message: "bank account created susscessfully",
		Data: map[string]interface{}{
			"bankAccount": newBankAccount,
		},
	}

	helper.WriteJson(w, http.StatusCreated, resp)

	return types.AppError{
		Error:  nil,
		Status: http.StatusCreated,
	}
}
