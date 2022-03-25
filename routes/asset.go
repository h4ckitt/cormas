package routes

import (
	"github.com/gofiber/fiber/v2"
	"rest-api/controllers"
)

func setupAssetRoutes(router fiber.Router) {
	assetGroup := router.Group("/asset")

	assetGroup.Post("/", authRequired(), controllers.CreateAsset)
	assetGroup.Delete("/:id", authRequired(), controllers.DeleteAsset)

}
