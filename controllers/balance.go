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

func CreateBalance(c *fiber.Ctx) error {
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

	balance := new(models.Balance)

	if err := c.BodyParser(&balance); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An Error Occurred While Processing That Request, Please Try Again Later",
		})
	}

	if balance.Status < 1 || balance.Status > 2 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Request	 Body Received",
		})
	}

	author := struct {
		UID string `json:"uid"`
	}{uid}

	balance.User = author
	balance.Type = "Balance"

	balanceJson, err := json.Marshal(balance)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "An Error Occurred While Processing That Request,, Please Try Again Later",
		})
	}

	mutation := &api.Mutation{
		CommitNow: true,
		SetJson:   balanceJson,
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

func GetBalance(c *fiber.Ctx) error {
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

	balanceID := c.Params("id")

	q := `
		query Balance($uid: string) {
			balance(func: uid($balanceUid)){
				amount
				currency {
					name
					value
					icon
					status
				}
				user(func: uid($uid)) {
					uid
					name
					email
					username
				}
				status
				moderation
				created_at
				updated_at
			}
		}
	`

	variables := map[string]string{
		"$uid": uid,
		"$balanceUid": balanceID,
	}

	resp, err := dgraph.NewTxn().QueryWithVars(context.Background(), q, variables)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An Error Occurred While Processing That Request",
		})
	}

	balance := struct {
		Result []models.Balance `json:"balance"`
	}{}

	_ = json.Unmarshal(resp.Json, &balance)

	if len(balance.Result) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Requested resource was not found on this server",
		})
	}

	return c.JSON(balance)
}