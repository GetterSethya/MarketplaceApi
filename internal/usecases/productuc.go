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

type ProductUseCase interface {
	CreateProduct(s datastore.Store, w http.ResponseWriter, r *http.Request) types.AppError
	UpdateProduct(s datastore.Store, w http.ResponseWriter, r *http.Request) types.AppError
	GetProductById(s datastore.Store, w http.ResponseWriter, r *http.Request) types.AppError
	ListProduct(s datastore.Store, w http.ResponseWriter, r *http.Request) types.AppError
	DeleteProduct(s datastore.Store, w http.ResponseWriter, r *http.Request) types.AppError
}

func UpdateStock(s datastore.Store, w http.ResponseWriter, r *http.Request) types.AppError {

	type stockStruct struct {
		Stock int `json:"stock"`
	}

	vars := mux.Vars(r)
	productIdUrlPath := vars["id"]
	userId := auth.GetUserIdFromJWT(r)

	body, err := io.ReadAll(r.Body)
	if err != nil {

		log.Println("error when reading body in productuc.go, updateStock()")

		return types.AppError{
			Error:  fmt.Errorf("Invalid/missing field"),
			Status: http.StatusOK,
		}
	}

	defer r.Body.Close()

	var stock *stockStruct

	err = json.Unmarshal(body, &stock)
	if err != nil {

		log.Println("error when Unmarshal body in update stock product usecase", err)

		return types.AppError{
			Error:  fmt.Errorf("Invalid/missing field"),
			Status: http.StatusBadRequest,
		}
	}

	product, err := s.GetProductById(productIdUrlPath)
	if err != nil {

		return types.AppError{
			Error:  fmt.Errorf("Product not found"),
			Status: http.StatusNotFound,
		}
	}

	if product.SellerId != userId {

		return types.AppError{
			Error:  fmt.Errorf("Forbidden"),
			Status: http.StatusForbidden,
		}
	}

	if err := s.UpdateStockProduct(productIdUrlPath, stock.Stock); err != nil {

		return types.AppError{
			Error:  fmt.Errorf("failed to update stock"),
			Status: http.StatusInternalServerError,
		}
	}

	resp := types.ServerResponse{
		Message: "Stock updated successfully",
		Data:    nil,
	}

	helper.WriteJson(w, http.StatusOK, resp)

	return types.AppError{
		Error:  nil,
		Status: http.StatusOK,
	}
}

func DeleteProduct(s datastore.Store, w http.ResponseWriter, r *http.Request) types.AppError {

	vars := mux.Vars(r)
	productIdUrlPath := vars["id"]
	userId := auth.GetUserIdFromJWT(r)
	sellerId, err := s.GetProductSeller(productIdUrlPath)
	if err != nil {
		log.Println("error when getting sellerid", err)

		return types.AppError{
			Error:  fmt.Errorf("Product didnot exist"),
			Status: http.StatusNotFound,
		}
	}

	if sellerId != userId {
		log.Println("error seller id != userid")

		return types.AppError{
			Error:  fmt.Errorf("Forbidden"),
			Status: http.StatusForbidden,
		}
	}

	if err := s.DeleteProduct(productIdUrlPath); err != nil {

		log.Println("error when deleting product", err)
		return types.AppError{
			Error:  fmt.Errorf("Product didnot exist"),
			Status: http.StatusNotFound,
		}
	}

	resp := types.ServerResponse{
		Message: "Product deleted successfully",
		Data:    nil,
	}

	helper.WriteJson(w, http.StatusOK, resp)

	return types.AppError{
		Error:  nil,
		Status: http.StatusOK,
	}
}

func ListProduct(s datastore.Store, w http.ResponseWriter, r *http.Request) types.AppError {
	//nampilin list product, GET /v1/product
	queries := getListProductQuery(r)
	userid := auth.GetUserIdFromJWT(r)

	//validasi query
	validQuery := validator.ValidateListProductQuery(queries)

	products, err := s.ListProducts(validQuery, userid)
	if err != nil {

		return types.AppError{
			Error:  fmt.Errorf("Error when fetching products"),
			Status: http.StatusInternalServerError,
		}
	}

	resp := types.ServerResponse{
		Message: "Ok",
		Data: map[string]interface{}{
			"products": products,
		},
	}

	helper.WriteJson(w, http.StatusOK, resp)

	return types.AppError{
		Error:  nil,
		Status: http.StatusOK,
	}
}

func GetProductById(s datastore.Store, w http.ResponseWriter, r *http.Request) types.AppError {

	vars := mux.Vars(r)
	productIdUrlPath := vars["id"]

	if !helper.ValidateUUID(productIdUrlPath) {

		return types.AppError{
			Error:  fmt.Errorf("Product didnot exist"),
			Status: http.StatusNotFound,
		}
	}

	product, err := s.GetProductById(productIdUrlPath)
	if err != nil {

		return types.AppError{
			Error:  fmt.Errorf("Product didnot exist"),
			Status: http.StatusNotFound,
		}
	}

	resp := types.ServerResponse{
		Message: "Ok",
		Data: map[string]interface{}{
			"product": product,
		},
	}

	helper.WriteJson(w, http.StatusOK, resp)

	return types.AppError{
		Error:  nil,
		Status: http.StatusOK,
	}
}

func UpdateProduct(s datastore.Store, w http.ResponseWriter, r *http.Request) types.AppError {

	vars := mux.Vars(r)
	productIdUrlPath := vars["id"]
	userId := auth.GetUserIdFromJWT(r)
	log.Println("userId: ", userId)

	productSellerId, err := s.GetProductSeller(productIdUrlPath)
	if err != nil {

		return types.AppError{
			Error:  fmt.Errorf("Product did not exist"),
			Status: http.StatusNotFound,
		}
	}

	if userId != productSellerId {

		return types.AppError{
			Error:  fmt.Errorf("Forbidden"),
			Status: http.StatusForbidden,
		}
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {

		log.Println("error when reading body from request in updateProduct in productuc.go")

		return types.AppError{
			Error:  fmt.Errorf("Invalid/missing field"),
			Status: http.StatusOK,
		}
	}

	defer r.Body.Close()

	var product *entities.Product

	err = json.Unmarshal(body, &product)
	if err != nil {

		log.Println("error when Unmarshal body in update product usecase", err)

		return types.AppError{
			Error:  fmt.Errorf("Invalid/missing field"),
			Status: http.StatusBadRequest,
		}
	}

	if err := validator.ValidateUpdateProductPayload(product); err != nil {

		log.Println("error when validating update product payload", err)

		return types.AppError{
			Error:  err,
			Status: http.StatusBadRequest,
		}
	}

	if err := s.UpdateProduct(productIdUrlPath, product); err != nil {

		log.Println("error when updating product in productuc.go:", err)

		return types.AppError{
			Error:  fmt.Errorf("Failed to update product"),
			Status: http.StatusInternalServerError,
		}
	}

	respProduct, err := s.GetProductById(productIdUrlPath)
	if err != nil {

		log.Println("error when getting product in productuc.go:", err)

		return types.AppError{
			Error:  fmt.Errorf("Failed to update product"),
			Status: http.StatusInternalServerError,
		}
	}

	resp := types.ServerResponse{
		Message: "Product updated susscessfully",
		Data: map[string]interface{}{
			"product": respProduct,
		},
	}

	helper.WriteJson(w, http.StatusOK, resp)

	return types.AppError{
		Error:  nil,
		Status: http.StatusOK,
	}
}

func CreateProduct(s datastore.Store, w http.ResponseWriter, r *http.Request) types.AppError {

	sellerId := auth.GetUserIdFromJWT(r)

	body, err := io.ReadAll(r.Body)
	if err != nil {

		log.Println("Error when reading body in CreateProduct usecase")

		return types.AppError{
			Error:  fmt.Errorf("Invalid/missing field"),
			Status: http.StatusBadRequest,
		}
	}

	defer r.Body.Close()

	var product *entities.Product

	err = json.Unmarshal(body, &product)
	if err != nil {

		log.Println("error when Unmarshal body in create product usecase", err)

		return types.AppError{
			Error:  fmt.Errorf("Invalid/missing field"),
			Status: http.StatusBadRequest,
		}
	}

	if err := validator.ValidateCreateProductPayload(product); err != nil {

		log.Println("error when validating create product payload")

		return types.AppError{
			Error:  err,
			Status: http.StatusBadRequest,
		}
	}

	id := uuid.NewString()

	if err := s.CreateProduct(id, sellerId, product); err != nil {

		log.Println("error when creating product", err)

		return types.AppError{
			Error:  fmt.Errorf("Failed when creating product, please try again."),
			Status: http.StatusInternalServerError,
		}
	}

	newProduct, err := s.GetProductById(id)

	if err != nil {

		log.Println("error when getting product", err)

		return types.AppError{
			Error:  fmt.Errorf("Failed when creating product, please try again."),
			Status: http.StatusInternalServerError,
		}
	}

	resp := types.ServerResponse{
		Message: "Product created susscessfully",
		Data: map[string]interface{}{
			"product": newProduct,
		},
	}

	helper.WriteJson(w, http.StatusCreated, resp)

	return types.AppError{
		Error:  nil,
		Status: http.StatusCreated,
	}
}

func getListProductQuery(r *http.Request) types.ListQuery {

	queryParams := r.URL.Query()

	// return product created by current user  "true"|"false"
	userOnly := queryParams.Get("useronly")

	// pagination
	limit := queryParams.Get("limit")
	offset := queryParams.Get("offset")

	// filter by tags
	tags := queryParams["tags"]

	// filter by conditions
	condition := queryParams.Get("condition")

	// return product by requester user that have 0 stock
	showEmptyStock := queryParams.Get("showemptystock")

	// return where product price bellow maxprice
	maxPrice := queryParams.Get("maxprice")

	// return where product price higher than minprice
	minPrice := queryParams.Get("minprice")

	// sort product by "asc"|"desc"
	sort := queryParams.Get("sort")

	// order product "price"|"date"|"name" default date
	order := queryParams.Get("order")

	// get product where name like search query
	search := queryParams.Get("search")

	return types.ListQuery{

		UserOnly:       userOnly,
		Limit:          limit,
		Offset:         offset,
		Tags:           tags,
		Condition:      condition,
		ShowEmptyStock: showEmptyStock,
		MaxPrice:       maxPrice,
		MinPrice:       minPrice,
		Sort:           sort,
		Order:          order,
		Search:         search,
	}

}
