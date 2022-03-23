package routes

import (
	"github.com/gofiber/fiber/v2"
	"rest-api/controllers"
)

func setupQuestionRoutes(router fiber.Router) {
	questionGroup := router.Group("/question")

	questionGroup.Post("/", authRequired(), controllers.CreateQuestion)
	questionGroup.Get("/", controllers.ListQuestions)
	questionGroup.Get("/:id", controllers.GetQuestion)
	questionGroup.Delete("/:id", authRequired(), controllers.DeleteQuestion)
	questionGroup.Patch("/:id", authRequired(), controllers.UpdateQuestion)
}
