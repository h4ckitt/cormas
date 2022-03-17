package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"rest-api/config"
	"rest-api/models"
	"time"
)

func GetJWTUser(token *jwt.Token) models.User {
	claims := token.Claims.(jwt.StandardClaims)
	uid := claims.Issuer

	//logic to fetch user from db using extracted uid
	fmt.Println(uid)

	return models.User{}
}

func GenerateJWTCookie(uid string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    uid,
		ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
	})

	fmt.Println(config.GetConfig().SecretKey)
	token, err := claims.SignedString([]byte(config.GetConfig().SecretKey))

	if err != nil {
		return "", err
	}

	return token, nil
}
