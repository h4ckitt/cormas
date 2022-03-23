package routes

import (
	"github.com/gofiber/fiber/v2"
	"rest-api/controllers"
)

func setupBankRoutes(router fiber.Router) {
	bankGroup := router.Group("/bank")

	bankGroup.Post("/", authRequired(), controllers.CreateBank)
	bankGroup.Get("/:id", authRequired(), controllers.GetBank)
	bankGroup.Delete("/:id", authRequired(), controllers.DeleteBank)
	bankGroup.Patch("/:id", authRequired(), controllers.UpdateBank)
}
