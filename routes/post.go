package routes

import (
	"github.com/gofiber/fiber/v2"
	"rest-api/controllers"
)

func setupPostRoutes(router fiber.Router) {
	postGroup := router.Group("/post")

	postGroup.Post("/", authRequired(), controllers.CreatePost)
	postGroup.Get("/:id", controllers.ReadPost)
	postGroup.Delete("/:id", authRequired(), controllers.DeletePost)
	postGroup.Patch("/:id", authRequired(), controllers.UpdatePost)
	postGroup.Post("/:id/comment", authRequired(), controllers.CreateComment)
	postGroup.Get("/:id/comment", controllers.ListComments)
}
