package routes

import (
	"github.com/gofiber/fiber/v2"
	"rest-api/controllers"
)

func setupCommentRoutes(router fiber.Router) {
	commentGroup := router.Group("/comment")

	commentGroup.Post("/", authRequired(), controllers.CreateComment)
	commentGroup.Delete("/:id", authRequired(), controllers.DeleteComment)
	commentGroup.Patch("/:id", authRequired(), controllers.UpdateComment)

}
