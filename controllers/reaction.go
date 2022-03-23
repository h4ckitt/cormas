package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"rest-api/models"
	"rest-api/utils"
)

// CreateReaction Discuss reaction model
func CreateReaction(c *fiber.Ctx) error {
	uid, err := utils.GetJWTUser(c.Locals("user").(*jwt.Token))

	if err != nil {
		if err.Error() == "invalid JWT Token" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Forbidden",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	_ := new(models.Reaction)
	return nil
}

func DeleteReaction(c *fiber.Ctx) error {

}
