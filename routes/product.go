package routes

import (
	"github.com/gofiber/fiber/v2"
	"rest-api/controllers"
)

func setupProductRoutes(router fiber.Router) {
	productGroup := router.Group("/product")

	productGroup.Post("/", authRequired(), controllers.CreateProduct)
	productGroup.Get("/", controllers.ListProducts)
	productGroup.Get("/:id", controllers.GetProduct)
	productGroup.Delete("/:id", authRequired(), controllers.DeleteProduct)
	productGroup.Patch("/:id", authRequired(), controllers.UpdateProduct)
	//productGroup.Post("/:id/review", authRequired(), controllers.Unimplemented)
	//productGroup.Get("/:id/comment", controllers.Unimplemented)
}
