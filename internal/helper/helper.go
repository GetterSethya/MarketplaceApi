package helper

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/GetterSethya/golangApiMarketplace/internal/types"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AppHandler func(w http.ResponseWriter, r *http.Request) types.AppError

func CreateHandlerFunc(f AppHandler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if err := f(w, r); err.Error != nil {
			resp := types.ServerResponse{
				Message: err.Error.Error(),
				Data:    nil,
			}

			WriteJson(w, err.Status, resp)
		}

	}
}

func WriteJson(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func GenerateHash(plainText string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainText), 12)
	if err != nil {
		log.Fatal("Error when hashing")
	}

	return string(hash)
}
func ArrayToString(arr []string) string {
	str := ""
	for i, v := range arr {
		if i > 0 {
			str += ","
		}
		str += v
	}
	return str
}

func ValidateUUID(id string) bool {

	_, err := uuid.Parse(id)

	return err == nil
}
