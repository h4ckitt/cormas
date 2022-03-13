package routes

import (
	"github.com/gofiber/fiber/v2"
)

func InitRouter() {
	app := fiber.New()

	v1 := app.Group("/v1", MiddleWare)

}
