package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"rest-api/models"
)

func GetJWTUser(token *jwt.Token) models.User {
	claims := token.Claims.(jwt.StandardClaims)
	uid := claims.Issuer

	//logic to fetch user from db using extracted uid
	fmt.Println(uid)

	return models.User{}
}
