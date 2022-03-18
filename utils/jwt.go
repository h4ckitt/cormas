package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"rest-api/db"
	"rest-api/models"
	"time"

	"rest-api/config"
)

var dgraph = db.GetDB()

func GetJWTUser(token *jwt.Token) (string, error) {
	claims := token.Claims.(jwt.MapClaims)

	uid := claims["uid"]
	q :=
		`
		query user($uid: string){
			user(func: uid($uid)) {
				suspended
			}
		}
		`

	resp, err := dgraph.NewTxn().QueryWithVars(context.Background(), q, map[string]string{"$uid": uid.(string)})

	if err != nil {
		return "", errors.New("error occurred while processing request")
	}

	user := struct {
		Users []models.User `json:"user"`
	}{}

	err = json.Unmarshal(resp.Json, &user)

	if err != nil {
		return "", errors.New("error marshalling json from db")
	}

	if len(user.Users) == 0 {
		return "", errors.New("invalid JWT Token")
	}

	return uid.(string), nil
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
