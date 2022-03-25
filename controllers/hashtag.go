package controllers

import (
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"log"
	"rest-api/models"
)

func ListTags(c *fiber.Ctx) error {
	q :=
		`
		{
			tags(func: type(HashTag)) {
				name
			}
		}
		`

	tags := struct {
		Result []models.HashTag `json:"tags"`
	}{}

	resp, err := dgraph.NewTxn().Query(context.Background(), q)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An Error Occurred While Processing That Request",
		})
	}

	json.Unmarshal(resp.Json, &tags)

	return c.JSON(tags.Result)
}
