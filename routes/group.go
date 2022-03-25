package routes

import (
	"github.com/gofiber/fiber/v2"
	"rest-api/controllers"
)

func setupGroupRoutes(router fiber.Router) {
	subGroup := router.Group("/group")

	subGroup.Post("/", authRequired(), controllers.CreateGroup)
	subGroup.Get("/all", authRequired(), controllers.ListGroups)
	subGroup.Delete("/:id", authRequired(), controllers.DeleteGroup)
	subGroup.Patch("/:id", authRequired(), controllers.UpdateGroup)
}
