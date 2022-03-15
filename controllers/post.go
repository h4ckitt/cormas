package controllers

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"rest-api/utils"
)

func CreatePost(c *fiber.Ctx) error {
	user := utils.GetJWTUser(c.Locals("user").(*jwt.Token))

	//	post := new(models.Post)

	//	mu := &api.Mutation{CommitNow: true}

	//create user posts
	fmt.Println(user)
	return nil
}

func ReadPost(c *fiber.Ctx) error {
	return nil
}

func DeletePost(c *fiber.Ctx) error {
	return nil
}

func UpdatePost(c *fiber.Ctx) error {
	return nil
}
