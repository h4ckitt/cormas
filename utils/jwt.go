package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"

	"rest-api/config"
)

func GetJWTUser(token *jwt.Token) string {
	claims := token.Claims.(jwt.MapClaims)

	uid := claims["uid"]

	return uid.(string)
}

func GenerateJWTCookie(uid string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":     uid,
		"expires": time.Now().Add(1 * time.Hour),
	})

	fmt.Println(config.GetConfig().SecretKey)
	token, err := claims.SignedString([]byte(config.GetConfig().SecretKey))

	if err != nil {
		return "", err
	}

	return token, nil
}
