package validator

import (
	"fmt"
	"strings"

	"github.com/GetterSethya/golangApiMarketplace/internal/entities"
)

const (
	MAXNAME     = 50
	MINNAME     = 5
	MAXUSERNAME = 15
	MINUSERNAME = 5
	MINPASSWORD = 5
	MAXPASSWORD = 15
)

func ValidateRegisterPayload(user *entities.User) error {

	var invalidFields []string

	nameLength := len(user.Name)
	usernameLength := len(user.Username)
	passwordLength := len(user.HashPassword)

	if user.Name == "" || nameLength < MINNAME || nameLength > MAXNAME {
		invalidFields = append(invalidFields, "name")
	}

	if user.Username == "" || usernameLength < MINUSERNAME || usernameLength > MAXUSERNAME || len(strings.Split(user.Username, " ")) > 1 {
		invalidFields = append(invalidFields, "username")
	}

	if user.HashPassword == "" || passwordLength < MINPASSWORD || passwordLength > MAXPASSWORD {
		invalidFields = append(invalidFields, "password")
	}

	if len(invalidFields) > 0 {
		return fmt.Errorf("Invalid " + strings.Join(invalidFields, ", "))
	}

	return nil
}

func ValidateLoginPayload(user *entities.User) error {

	var invalidFields []string

	usernameLength := len(user.Username)
	passwordLength := len(user.HashPassword)

	if user.Username == "" || usernameLength < MINUSERNAME || usernameLength > MAXUSERNAME || len(strings.Split(user.Username, " ")) > 1 {
		invalidFields = append(invalidFields, "username")
	}

	if user.HashPassword == "" || passwordLength < MINPASSWORD || passwordLength > MAXPASSWORD {
		invalidFields = append(invalidFields, "password")
	}

	if len(invalidFields) > 0 {
		return fmt.Errorf("Invalid " + strings.Join(invalidFields, ", "))
	}

	return nil
}

func ValidateUpdatePayload(user *entities.User) error {

	var invalidFields []string

	nameLength := len(user.Name)
	usernameLength := len(user.Username)

	if user.Name == "" || nameLength < MINNAME || nameLength > MAXNAME {
		invalidFields = append(invalidFields, "name")
	}

	if user.Username == "" || usernameLength < MINUSERNAME || usernameLength > MAXUSERNAME || len(strings.Split(user.Username, " ")) > 1 {
		invalidFields = append(invalidFields, "username")
	}

	if len(invalidFields) > 0 {
		return fmt.Errorf("Invalid " + strings.Join(invalidFields, ", "))
	}

	return nil
}
