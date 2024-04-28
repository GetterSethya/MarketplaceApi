package auth

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/GetterSethya/golangApiMarketplace/config"
	"github.com/GetterSethya/golangApiMarketplace/internal/helper"
	"github.com/GetterSethya/golangApiMarketplace/internal/types"
	"github.com/golang-jwt/jwt/v4"
)

// middleware untuk mem-protect route, jika tidak ada Header "Authorization" atau JWT tidak valid, maka akan mereturn error status forbidden 403
func JWTMiddleware(f helper.AppHandler) helper.AppHandler {

	return func(w http.ResponseWriter, r *http.Request) types.AppError {
		secret := config.LoadConfig().App.JWTSecret

		// get Authorization header
		jwtToken := getTokenFromRequest(r)

		//
		token, err := validateJWT(jwtToken, secret)
		if err != nil || !token.Valid {

			return types.AppError{
				Error:  fmt.Errorf("Invalid token"),
				Status: http.StatusForbidden,
			}
		}

		// call appHandler func
		if err := f(w, r); err.Error != nil {

			return err
		}

		return types.AppError{
			Error:  nil,
			Status: http.StatusOK,
		}
	}

}

// subject berisi userId
//
// issuer "shopifyx"
// expiration time default 12 hours, as time after the token issued + 12 hours
//
// not before: The "nbf" (not before) claim identifies the time before which the JWT
//
// issued At: The "iat" (issued at) claim identifies the time at which the JWT was issued.  This claim can be used to determine the age of the JWT.MUST NOT be accepted for processing
func CreateJWT(userId, secret string) (string, error) {
	exp := jwt.NewNumericDate(time.Now().Add(time.Hour * 12))
	nbf := jwt.NewNumericDate(time.Now())
	iat := jwt.NewNumericDate(time.Now())
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userId,
		Issuer:    "shopifyx",
		ExpiresAt: exp,
		NotBefore: nbf,
		IssuedAt:  iat,
	})

	accessToken, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Printf("Error when signing accessToken %+v", err.Error())
		return "", err
	}

	return accessToken, nil
}

// hanya panggil fungsi ini di route yang sudah ada dijaga oleh middleware JWTMiddleware
func GetUserIdFromJWT(r *http.Request) string {
	secret := config.LoadConfig().App.JWTSecret
	token := getTokenFromRequest(r)
	jwtToken, err := validateJWT(token, secret)

	if err != nil {
		return ""
	}

	return jwtToken.Claims.(jwt.MapClaims)["sub"].(string)
}

func validateJWT(token, secret string) (*jwt.Token, error) {

	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %+v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
}

func getTokenFromRequest(r *http.Request) string {
	jwtToken := r.Header.Get("Authorization")

	return jwtToken
}
