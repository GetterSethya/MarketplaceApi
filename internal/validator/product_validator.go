package validator

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/GetterSethya/golangApiMarketplace/internal/entities"
	"github.com/GetterSethya/golangApiMarketplace/internal/types"
)

const (
	MAXPRODUCTNAME = 200
	MINPRODUCTNAME = 2
	MAXPRICE       = 10000000
	MINPRICE       = 0
	MAXIMAGEURL    = 255
	MAXSTOCK       = 32000
)

func ValidateListProductQuery(q types.ListQuery) types.ListQueryValid {
	userOnly := strings.ToLower(q.UserOnly)
	limit, err := strconv.Atoi(q.Limit)
	if err != nil || limit < 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(q.Offset)
	if err != nil || offset < 0 {
		offset = 0
	}

	condition := strings.ToLower(q.Condition)
	showEmptyStock := strings.ToLower(q.ShowEmptyStock)
	sort := strings.ToLower(q.Sort)
	order := strings.ToLower(q.Order)

	minPrice, err := strconv.Atoi(q.MinPrice)
	if err != nil {
		minPrice = 0
	}

	maxPrice, err := strconv.Atoi(q.MaxPrice)
	if err != nil {
		maxPrice = 0
	}

	if !(userOnly == "false" || userOnly == "true") {
		userOnly = "false"
	}

	if !(condition == "new" || condition == "second") {
		condition = "new"
	}

	if !(showEmptyStock == "false" || showEmptyStock == "true") {
		showEmptyStock = "false"
	}

	if !(sort == "asc" || sort == "desc") {
		sort = "desc"
	}

	if order == "date" {
		order = "createdAt"
	}

	if !(order == "createdAt" || order == "name" || sort == "price") {
		order = "createdAt"
	}

	return types.ListQueryValid{
		UserOnly:       userOnly,
		Limit:          limit,
		Offset:         offset,
		Tags:           q.Tags,
		Condition:      condition,
		ShowEmptyStock: showEmptyStock,
		MaxPrice:       float64(maxPrice),
		MinPrice:       float64(minPrice),
		Sort:           sort,
		Order:          order,
		Search:         q.Search,
	}
}

func ValidateUpdateProductPayload(p *entities.Product) error {

	var invalidFields []string
	nameLength := len(p.Name)

	if p.Name == "" || nameLength < MINPRODUCTNAME || nameLength > MAXPRODUCTNAME {
		invalidFields = append(invalidFields, "product name")
	}

	if p.Price < MINPRICE || p.Price > MAXPRICE {
		invalidFields = append(invalidFields, "product price")
	}

	if p.ImageUrl == "" || len(p.ImageUrl) > MAXIMAGEURL {
		invalidFields = append(invalidFields, "product imageUrl")
	}

	if p.Stock < 0 || p.Stock > MAXSTOCK {
		invalidFields = append(invalidFields, "product stock")
	}

	if !validateCondition(p.Condition) {
		invalidFields = append(invalidFields, "product condition")
	}

	if len(invalidFields) > 0 {
		return fmt.Errorf("Invalid " + strings.Join(invalidFields, ", "))
	}

	return nil
}

func ValidateCreateProductPayload(p *entities.Product) error {

	var invalidFields []string
	nameLength := len(p.Name)

	if p.Name == "" || nameLength < MINPRODUCTNAME || nameLength > MAXPRODUCTNAME {
		invalidFields = append(invalidFields, "product name")
	}

	if p.Price < MINPRICE || p.Price > MAXPRICE {
		invalidFields = append(invalidFields, "product price")
	}

	if p.ImageUrl == "" || len(p.ImageUrl) > MAXIMAGEURL {
		invalidFields = append(invalidFields, "product imageUrl")
	}

	if p.Stock < 0 || p.Stock > MAXSTOCK {
		invalidFields = append(invalidFields, "product stock")
	}

	if !validateCondition(p.Condition) {
		invalidFields = append(invalidFields, "product condition")
	}

	if len(invalidFields) > 0 {
		return fmt.Errorf("Invalid " + strings.Join(invalidFields, ", "))
	}

	return nil
}

func validateCondition(condition string) bool {
	return condition == "new" || condition == "second"
}
