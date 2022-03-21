package controllers

import (
	"context"
	"encoding/json"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"rest-api/models"
	"rest-api/utils"
)

func CreatePost(c *fiber.Ctx) error {
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

	post := new(models.Post)

	if err := c.BodyParser(&post); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An Error Occurred While Processing That Request, Please Try Again Later",
		})
	}

	if post.Name == "" || post.Description == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Request	 Body Received",
		})
	}

	author := struct {
		UID string `json:"uid"`
	}{uid}

	post.Author = author
	post.Type = "Post"

	postJson, err := json.Marshal(post)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "An Error Occurred While Processing That Request,, Please Try Again Later",
		})
	}

	mutation := &api.Mutation{
		CommitNow: true,
		SetJson:   postJson,
	}

	_, err = dgraph.NewTxn().Mutate(context.Background(), mutation)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "An Error Occurred While Processing That Request,, Please Try Again Later",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Post Created Successfully",
	})
}

func ReadPost(c *fiber.Ctx) error {
	//	uid := c.Params("id")

	return nil
}

func DeletePost(c *fiber.Ctx) error {
	return nil
}

func UpdatePost(c *fiber.Ctx) error {
	return nil
}
