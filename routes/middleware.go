package routes

import (
	"fmt"

	"github.com/gofiber/fiber"
)

func MiddleWare(c *fiber.Ctx) error {
	fmt.Println("Request")
}
