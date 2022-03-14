package routes

import (
	"github.com/gofiber/fiber/v2"
)

func InitRouter() *fiber.App {
	app := fiber.New()

	v1 := app.Group("/v1")

	setupUserRoutes(v1)
	setupAnswerRoutes(v1)
	setupPostRoutes(v1)

	return app
}
