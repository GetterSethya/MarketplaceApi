package validator

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/GetterSethya/golangApiMarketplace/internal/entities"
	"github.com/GetterSethya/golangApiMarketplace/internal/helper"
	"github.com/GetterSethya/golangApiMarketplace/internal/types"
)

const (
	MAXQTT         = 32000
	MAXNOTESLENGTH = 1
	MINQTT         = 1
	MINNOTESLENGTH = 1
)

func ValidateListTransactionQuery(q types.ListQueryTransaction) types.ListQueryTransactionValid {

	seller := strings.ToLower(q.Seller)
	sort := strings.ToLower(q.Sort)
	order := strings.ToLower(q.Order)

	if !(seller == "true" || seller == "false") {
		seller = "false"
	}

	parsedSeller, err := strconv.ParseBool(seller)
	if err != nil {
		parsedSeller = false
	}

	limit, err := strconv.Atoi(q.Limit)
	if err != nil || limit < 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(q.Offset)
	if err != nil || offset < 0 {
		offset = 0
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

	return types.ListQueryTransactionValid{
		Seller: parsedSeller,
		Limit:  limit,
		Offset: offset,
		Sort:   sort,
		Order:  order,
		Search: q.Search,
	}
}

func ValidateUpdateStatusTransactionPayload(status string) error {

	if !(status == "menunggu" || status == "diterima seller" || status == "dalam pengiriman" || status == "diterima" || status == "ditolak") {
		return fmt.Errorf("Invalid status")
	}

	return nil
}

func ValidateCreateTransactionPayload(p *entities.Transaction) error {

	var invalidFields []string
	notesLength := len(p.Notes)

	if ok := helper.ValidateUUID(p.ProductId); ok == false {
		invalidFields = append(invalidFields, "transaction productId")
	}

	if notesLength > MAXNOTESLENGTH {
		invalidFields = append(invalidFields, "transaction notes")
	}

	if p.Quantity < MINQTT || p.Quantity > MAXQTT {
		invalidFields = append(invalidFields, "transaction quantity")
	}

	if len(invalidFields) > 0 {
		return fmt.Errorf("Invalid " + strings.Join(invalidFields, ", "))
	}

	return nil
}
