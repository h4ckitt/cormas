package routes

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func InitRouter() *fiber.App {

	fmt.Println("Here")
	app := fiber.New()

	v1 := app.Group("/v1")

	//setupAnswerRoutes(v1)
	setupUserRoutes(v1)
	setupPostRoutes(v1)
	setupCurrencyRoutes(v1)
	setupCommentRoutes(v1)
	setupProductRoutes(v1)

	return app
}
