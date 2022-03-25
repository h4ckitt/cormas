package routes

import (
	"github.com/gofiber/fiber/v2"
	"rest-api/controllers"
)

func setupTagRoutes(router fiber.Router) {
	tagGroup := router.Group("/hashtag")

	tagGroup.Get("/", controllers.ListTags)
	//tagGroup.Delete("/:id", authRequired(), controllers.DeleteBank) //Only An Admin Should Be Able To Do This?
}
