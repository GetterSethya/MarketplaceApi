package validator

import (
	"fmt"
	"strings"

	"github.com/GetterSethya/golangApiMarketplace/internal/entities"
)

const (
	MAXBANKNAME    = 50
	MAXACCOUNTNAME = 100
)

func ValidateUpdateBankAccountPayload(b *entities.BankAccount) error {
	var invalidFields []string
	bankNameLength := len(b.BankName)
	accountNameLength := len(b.AccountName)

	if b.BankName == "" || bankNameLength < 1 || bankNameLength > MAXBANKNAME {
		invalidFields = append(invalidFields, "bank name")
	}

	if b.AccountName == "" || accountNameLength < 1 || accountNameLength > MAXACCOUNTNAME {
		invalidFields = append(invalidFields, "account name")
	}

	if b.AccountNumber == 0 {
		invalidFields = append(invalidFields, "account number")
	}

	if len(invalidFields) > 0 {
		return fmt.Errorf("Invalid " + strings.Join(invalidFields, ", "))
	}

	return nil
}

func ValidateCreateBankAccountPayload(b *entities.BankAccount) error {

	var invalidFields []string
	bankNameLength := len(b.BankName)
	accountNameLength := len(b.AccountName)

	if b.BankName == "" || bankNameLength < 1 || bankNameLength > MAXBANKNAME {
		invalidFields = append(invalidFields, "bank name")
	}

	if b.AccountName == "" || accountNameLength < 1 || accountNameLength > MAXACCOUNTNAME {
		invalidFields = append(invalidFields, "account name")
	}

	if b.AccountNumber == 0 {
		invalidFields = append(invalidFields, "account number")
	}

	if len(invalidFields) > 0 {
		return fmt.Errorf("Invalid " + strings.Join(invalidFields, ", "))
	}

	return nil
}
