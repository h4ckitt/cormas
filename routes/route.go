package routes

import (
	"github.com/gofiber/fiber/v2"
	"rest-api/controllers"
)

//const secretkey = "Mr.RavandraIsTheB3sTUpw0rKCl13nt"
func InitRouter() *fiber.App {
	app := fiber.New()

	/*app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(secretkey),
	}))*/
	//x := controllers.SignUpHandler
	v1 := app.Group("/v1")

	userGroup := v1.Group("/user")
	userGroup.Get("/", controllers.Unimplemented)
	return app
}
