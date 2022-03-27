package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"rest-api/models"
	"rest-api/utils"
)

func CreateCategory(c *fiber.Ctx) error {
	uid, err := utils.GetJWTUser(c.Locals("user").(*jwt.Token))

	if err != nil {
		if err.Error() == "invalid JWT Token" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Forbidden",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	category := new(models.Category)

	if err := c.BodyParser(category); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request body received",
		})
	}

	if category.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request body received",
		})
	}

	category.Type = "Category"

	//categoryJson, err := json.Marshal(category)
	return nil
}

func ListCategories(c *fiber.Ctx) error {
	return nil
}

func DeleteCategory(c *fiber.Ctx) error {
	return nil
}
func CreateSubCategory(c *fiber.Ctx) error {
	return nil
}
