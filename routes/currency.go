package routes

import (
	"github.com/gofiber/fiber/v2"
	"rest-api/controllers"
)

func setupCurrencyRoutes(router fiber.Router) {
	currencyGroup := router.Group("/currency")

	currencyGroup.Post("", controllers.CreateCurrency)
	currencyGroup.Get("/all", controllers.GetAllCurrencies)
	currencyGroup.Get("/:uid", controllers.GetOneCurrency)
	currencyGroup.Put("/:uid", controllers.UpdateCurrency)
	currencyGroup.Delete("/:uid", controllers.DeleteCurrency)
}
