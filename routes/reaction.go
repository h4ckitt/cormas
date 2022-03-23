package routes

import (
	"github.com/gofiber/fiber/v2"
	"rest-api/controllers"
)

func setupReactionRoutes(router fiber.Router) {
	reactionGroup := router.Group("/reaction")

	reactionGroup.Post("/", authRequired(), controllers.CreateReaction)
	reactionGroup.Delete("/", authRequired(), controllers.DeleteReaction)
}
