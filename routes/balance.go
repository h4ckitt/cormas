package routes

import (
	"github.com/gofiber/fiber/v2"
	"rest-api/controllers"
)

func setupBalanceRoutes(router fiber.Router) {
	balanceGroup := router.Group("/balance")

	balanceGroup.Post("/", authRequired(), controllers.CreateBalance)
	balanceGroup.Get("/:id", authRequired(), controllers.GetBalance)
	balanceGroup.Delete("/:id", authRequired(), controllers.DeleteBalance)
	balanceGroup.Patch("/:id", authRequired(), controllers.UpdateBalance)

}