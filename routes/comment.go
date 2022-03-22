package routes

import (
	"github.com/gofiber/fiber/v2"
	"rest-api/controllers"
)

func setupCommentRoutes(router fiber.Router) {
	commentGroup := router.Group("/comment")

	//commentGroup.Post("/", authRequired(), controllers.CreatePost)
	//commentGroup.Get("/:id", controllers.ReadPost)
	commentGroup.Delete("/:id", authRequired(), controllers.DeleteComment)
	commentGroup.Patch("/:id", authRequired(), controllers.UpdateComment)
	//commentGroup.Post("/:id/comment", authRequired(), controllers.CreateComment)
}
