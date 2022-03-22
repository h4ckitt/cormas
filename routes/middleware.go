package routes

import (
	"fmt"
	jwtware "github.com/gofiber/jwt/v3"
	"rest-api/config"

	"github.com/gofiber/fiber/v2"
)

/*var JWTMiddleWare = jwtware.New(jwtware.Config{
	SigningKey: []byte(config.GetConfig().SecretKey),
})*/

func authRequired() func(c *fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Forbidden",
			})
		},

		SigningKey: []byte(config.GetConfig().SecretKey),
	})
}

func MiddleWare(c *fiber.Ctx) error {
	fmt.Println("Request")
	return nil
}
