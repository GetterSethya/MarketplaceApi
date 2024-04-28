package auth

import (
	"log"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func TestValidateJWT(t *testing.T) {
	userId := "12345678"
	secret := "superSecret"

	jwtString, err := CreateJWT(userId, secret)
	if err != nil {
		log.Println(jwtString)
		t.Errorf("Failed when creating jwt")
	}

	validatedToken, err := validateJWT(jwtString, secret)
	if err != nil {
		t.Errorf("Failed when validating jwt token")
	}

	sub := validatedToken.Claims.(jwt.MapClaims)["sub"].(string)
	if sub == "" {
		t.Errorf("Expected %s. got=%s", userId, sub)
	}

	if sub != userId {
		t.Errorf("Expected userId to be %s, but got=%s", userId, sub)
	}

	exp := validatedToken.Claims.(jwt.MapClaims)["exp"].(float64)
	if int64(exp) < time.Now().Unix() {
		t.Errorf("Invalid exp, expected exp bigger than current time. but got=%+v", time.Unix(int64(exp), 0))
	}

	if err := validatedToken.Claims.Valid(); err != nil {
		t.Errorf("Failed to validate claims")
	}
}
