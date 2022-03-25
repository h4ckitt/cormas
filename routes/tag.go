package routes

import (
	"github.com/gofiber/fiber/v2"
	"rest-api/controllers"
)

func setupTagRoutes(router fiber.Router) {
	tagGroup := router.Group("/tag")

	tagGroup.Get("/:id", controllers.ListTags)
	//tagGroup.Delete("/:id", authRequired(), controllers.DeleteBank) //Only An Admin Should Be Able To Do This?
}
