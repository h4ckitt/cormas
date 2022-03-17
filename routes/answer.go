package routes

import (
	"github.com/gofiber/fiber/v2"
	"rest-api/controllers"
)

func setupAnswerRoutes(router fiber.Router) {
	answerGroup := router.Group("/answer")
	answerGroup.Post("/", controllers.CreateAnswer)
}
