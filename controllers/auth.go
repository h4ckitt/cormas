package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"rest-api/config"
	"rest-api/db"
	"rest-api/models"
	"strconv"
	"time"
)

func SignUpHandler(c *fiber.Ctx) error {
	dgraph := db.GetDB()

	txn := dgraph.NewTxn()

	defer func(txn *dgo.Txn, ctx context.Context) {
		err := txn.Discard(ctx)
		if err != nil {

		}
	}(txn, context.TODO())

	user := new(models.User)
	//var data map[string]string

	if err := c.BodyParser(user); err != nil {
		fmt.Println("An Error Occured: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	if user.Name == "" || user.Password == "" || user.Email == "" || user.Username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "bad request body received",
		})
	}

	//fmt.Println(data)

	password, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 12)

	user.Password = string(password)

	userBody, err := json.Marshal(user)

	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "An Error Occurred, Please Try Again",
		})
	}

	mu := &api.Mutation{
		SetJson: userBody,
	}

	mu.CommitNow = true

	_, err = txn.Mutate(context.TODO(), mu)

	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "An Error Occurred",
		})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func Login(c *fiber.Ctx) error {
	var (
		data map[string]string
		user models.User
	)

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	//fetch user info from db logic here

	if user.UID == 0 {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Invalid Credentials",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"])); err != nil {
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid Credentials",
		})
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.UID)),
		ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
	})

	token, err := claims.SignedString([]byte(config.GetConfig().SecretKey))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An Error Occurred, Please Try Again Later",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "authorization_token",
		Value:    token,
		Expires:  time.Now().Add(1 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User Logged In Successfully",
	})
}
