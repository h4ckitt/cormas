package routes

import (
	"github.com/gofiber/fiber/v2"
	"rest-api/controllers"
)

func setupReviewRoutes(router fiber.Router) {
	reviewGroup := router.Group("/review")

	reviewGroup.Post("/", authRequired(), controllers.CreateReview)
	reviewGroup.Patch("/:id", authRequired(), controllers.UpdateReview)
	reviewGroup.Delete("/:id", authRequired(), controllers.DeleteReview)
}
