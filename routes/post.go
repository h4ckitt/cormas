package routes

import (
	"github.com/gofiber/fiber/v2"
	"rest-api/controllers"
)

func setupPostRoutes(router fiber.Router) {
	postGroup := router.Group("/post")

	postGroup.Post("/", authRequired(), controllers.CreatePost)
}
