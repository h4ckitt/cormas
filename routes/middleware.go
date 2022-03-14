package routes

import (
	"fmt"
	jwtware "github.com/gofiber/jwt/v3"
	"rest-api/config"

	"github.com/gofiber/fiber/v2"
)

var JWTMiddleWare = jwtware.New(jwtware.Config{
	SigningKey: []byte(config.GetConfig().SecretKey),
})

func MiddleWare(c *fiber.Ctx) error {
	fmt.Println("Request")
	return nil
}
