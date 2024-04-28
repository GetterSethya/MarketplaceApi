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

type TransactionUseCase interface {
	CreateTransaction(s datastore.Store, w http.ResponseWriter, r *http.Request) types.AppError
	GetTransaction(s datastore.Store, w http.ResponseWriter, r *http.Request) types.AppError
	ListTransaction(s datastore.Store, w http.ResponseWriter, r *http.Request) types.AppError
	UpdateStatusTransaction(s datastore.Store, w http.ResponseWriter, r *http.Request) types.AppError
}

func CreateTransaction(s datastore.Store, w http.ResponseWriter, r *http.Request) types.AppError {

	buyerId := auth.GetUserIdFromJWT(r)
	body, err := io.ReadAll(r.Body)
	if err != nil {

		log.Println("Error when reading body in CreateTransaction usecases")

		return types.AppError{
			Error:  fmt.Errorf("Invalid/missing field"),
			Status: http.StatusBadRequest,
		}
	}

	defer r.Body.Close()

	var transaction *entities.Transaction

	err = json.Unmarshal(body, &transaction)
	if err != nil {

		log.Println("error when Unmarshal body in create transaction usecases", err)

		return types.AppError{
			Error:  fmt.Errorf("Invalid/missing field"),
			Status: http.StatusBadRequest,
		}
	}

	if err := validator.ValidateCreateTransactionPayload(transaction); err != nil {
		log.Println("error when validating create transaction payload")

		return types.AppError{
			Error:  err,
			Status: http.StatusBadRequest,
		}
	}

	product, err := s.GetProductById(transaction.ProductId)
	if err != nil {
		log.Println("error when creating transaction", err)

		return types.AppError{
			Error:  fmt.Errorf("Failed when creating transaction, product didnot exist"),
			Status: http.StatusBadRequest,
		}
	}

	if buyerId == product.SellerId {

		return types.AppError{
			Error:  fmt.Errorf("Cannot buy your own product"),
			Status: http.StatusBadRequest,
		}
	}

	total := product.Price * float64(transaction.Quantity)
	id := uuid.NewString()

	if err := s.CreateTransaction(id, buyerId, product.SellerId, product.ID, total, transaction); err != nil {

		log.Println("error when creating transaction", err)

		return types.AppError{
			Error:  fmt.Errorf("Failed when creating transaction, please try again."),
			Status: http.StatusInternalServerError,
		}
	}

	newTransaction, err := s.GetTransaction(id)
	if err != nil {

		log.Println("error when getting transaction", err)

		return types.AppError{
			Error:  fmt.Errorf("Failed when creating transaction, please try again."),
			Status: http.StatusInternalServerError,
		}
	}

	if err := s.UpdateStockProduct(product.ID, product.Stock-1); err != nil {
		log.Println("error when updating product stock in create transaction", err)
		return types.AppError{
			Error:  fmt.Errorf("Failed when creating transaction, please try again."),
			Status: http.StatusInternalServerError,
		}
	}

	resp := types.ServerResponse{
		Message: "Transaction created susscessfully",
		Data:    newTransaction,
	}

	helper.WriteJson(w, http.StatusCreated, resp)

	return types.AppError{
		Error:  nil,
		Status: http.StatusCreated,
	}
}

func GetTransaction(s datastore.Store, w http.ResponseWriter, r *http.Request) types.AppError {

	vars := mux.Vars(r)
	transactionIdUrlPath := vars["id"]
	userId := auth.GetUserIdFromJWT(r)

	println("userid", userId)

	if !helper.ValidateUUID(transactionIdUrlPath) {

		return types.AppError{
			Error:  fmt.Errorf("Product didnot exist"),
			Status: http.StatusNotFound,
		}
	}

	transaction, err := s.GetTransaction(transactionIdUrlPath)
	if err != nil {

		return types.AppError{
			Error:  fmt.Errorf("Transaction didnot exist"),
			Status: http.StatusNotFound,
		}
	}
	println("transaction.buyer.id", transaction.Buyer.ID)
	println("transaction.seller.id", transaction.Seller.ID)

	if !(transaction.Buyer.ID == userId || transaction.Seller.ID == userId) {

		return types.AppError{
			Error:  fmt.Errorf("Forbidden"),
			Status: http.StatusForbidden,
		}
	}

	resp := types.ServerResponse{
		Message: "Ok",
		Data: map[string]interface{}{
			"transaction": transaction,
		},
	}

	helper.WriteJson(w, http.StatusOK, resp)

	return types.AppError{
		Error:  nil,
		Status: http.StatusOK,
	}
}

func ListTransaction(s datastore.Store, w http.ResponseWriter, r *http.Request) types.AppError {
	queries := getListTransactionQuery(r)
	userId := auth.GetUserIdFromJWT(r)
	validQuery := validator.ValidateListTransactionQuery(queries)

	transactions, err := s.ListTransaction(validQuery, userId)
	if err != nil {

		return types.AppError{
			Error:  fmt.Errorf("Failed when fetching transactions"),
			Status: http.StatusInternalServerError,
		}
	}

	resp := types.ServerResponse{
		Message: "Ok",
		Data: map[string]interface{}{
			"transactions": transactions,
		},
	}

	helper.WriteJson(w, http.StatusOK, resp)

	return types.AppError{
		Error:  nil,
		Status: http.StatusOK,
	}
}

func UpdateStatusTransaction(s datastore.Store, w http.ResponseWriter, r *http.Request) types.AppError {

	vars := mux.Vars(r)
	transactionIdUrlPath := vars["id"]
	userId := auth.GetUserIdFromJWT(r)

	if !helper.ValidateUUID(transactionIdUrlPath) {

		return types.AppError{
			Error:  fmt.Errorf("Product didnot exist"),
			Status: http.StatusNotFound,
		}
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {

		log.Println("Error when reading body in UpdateStatusTransaction usecases")

		return types.AppError{
			Error:  fmt.Errorf("Invalid/missing field"),
			Status: http.StatusBadRequest,
		}
	}

	defer r.Body.Close()

	var transaction *entities.Transaction

	err = json.Unmarshal(body, &transaction)
	if err != nil {

		log.Println("error when Unmarshal body in create transaction usecases", err)

		return types.AppError{
			Error:  fmt.Errorf("Invalid/missing field"),
			Status: http.StatusBadRequest,
		}
	}

	if err := validator.ValidateUpdateStatusTransactionPayload(transaction.Status); err != nil {

		log.Println("error when validating create transaction payload")

		return types.AppError{
			Error:  err,
			Status: http.StatusBadRequest,
		}
	}

	tx, err := s.GetTransaction(transactionIdUrlPath)
	if err != nil {

		return types.AppError{
			Error:  fmt.Errorf("Transaction didnot exist"),
			Status: http.StatusNotFound,
		}
	}

	if tx.Buyer.ID != userId || tx.Seller.ID != userId {

		return types.AppError{
			Error:  fmt.Errorf("Forbidden"),
			Status: http.StatusForbidden,
		}
	}

	if tx.Buyer.ID == userId && !(transaction.Status == "diterima") {
		return types.AppError{
			Error:  fmt.Errorf("Forbidden"),
			Status: http.StatusForbidden,
		}
	}

	if tx.Seller.ID == userId && !(transaction.Status == "diterima seller" || transaction.Status == "dalam pengiriman") {
		return types.AppError{
			Error:  fmt.Errorf("Forbidden"),
			Status: http.StatusForbidden,
		}
	}

	if err := s.UpdateStatusTransaction(tx.Transaction.ID, transaction.Status); err != nil {

		return types.AppError{
			Error:  fmt.Errorf("Failed updating transaction, something went wrong"),
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

func getListTransactionQuery(r *http.Request) types.ListQueryTransaction {

	queryParams := r.URL.Query()

	// return transaction list as seller
	seller := queryParams.Get("seller")

	// pagination
	limit := queryParams.Get("limit")
	offset := queryParams.Get("offset")

	// sort product by "asc"|"desc"
	sort := queryParams.Get("sort")

	// order product "price"|"date"|"name" default date
	order := queryParams.Get("order")

	// get product where id == productId
	search := queryParams.Get("search")

	return types.ListQueryTransaction{
		Seller: seller,
		Limit:  limit,
		Offset: offset,
		Sort:   sort,
		Order:  order,
		Search: search,
	}
}
