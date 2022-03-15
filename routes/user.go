package routes

import (
	"github.com/gofiber/fiber/v2"
	"rest-api/controllers"
)

func setupUserRoutes(router fiber.Router) {
	userGroup := router.Group("/user")
	userGroup.Post("/", controllers.SignUpHandler)
	userGroup.Post("/login", controllers.Login)

	//userGroup.Use(JWTMiddleWare)

}
