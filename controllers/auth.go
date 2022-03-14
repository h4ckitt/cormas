package controllers

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"rest-api/config"
	"rest-api/models"
	"strconv"
	"time"
)

func decodeToken(token *jwt.Token) models.User {
	claims := token.Claims.(jwt.StandardClaims)
	uid := claims.Issuer

	//logic to fetch user from db using extracted uid
	fmt.Println(uid)

	return models.User{}
}

func SignUpHandler(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 12)

	user := models.User{
		Name:     data["name"],
		Email:    data["email"],
		Password: string(password),
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
