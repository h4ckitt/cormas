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
	"time"
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
	now := time.Now().Format(time.RFC3339)

	balance.User = author
	balance.CreatedAt = now
	balance.UpdatedAt = now
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
				name
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

func DeleteBalance(c *fiber.Ctx) error {
	uid, _ := utils.GetJWTUser(c.Locals("user").(*jwt.Token))

	balanceId := c.Params("id")

	q :=
		`
		query Balance($uid: string) {
			balance(func: uid($uid)) @normalize{
				name: name
				User {
					user_uid: uid
				}
			}
		}
		`

	resp, err := dgraph.NewTxn().QueryWithVars(context.Background(), q, map[string]string{"$uid": balanceId})

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	balance := struct {
		Result []struct {
			Name      string `json:"name"`
			UserUid string `json:"user_uid"`
		} `json:"balance"`
	}{}

	_ = json.Unmarshal(resp.Json, &balance)

	if len(balance.Result) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "The resource doesn't exist on this server",
		})
	}

	if balance.Result[0].UserUid != uid {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
		})
	}

	tbd := struct {
		UID string `json:"uid"`
	}{balanceId}

	tbdJson, err := json.Marshal(tbd)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	mutation := &api.Mutation{
		CommitNow:  true,
		DeleteJson: tbdJson,
	}

	_, err = dgraph.NewTxn().Mutate(context.Background(), mutation)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Deleted",
	})
}

func UpdateBalance(c *fiber.Ctx) error {
	uid, _ := utils.GetJWTUser(c.Locals("user").(*jwt.Token))

	balanceId := c.Params("id")

	if balanceId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request received",
		})
	}

	tbu := struct {
		UID         string `json:"uid"`
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`
		UpdatedAt   string `json:"updated_at"`
	}{UID: balanceId, UpdatedAt: time.Now().Format(time.RFC3339)}

	if err := c.BodyParser(&tbu); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	if tbu.Description == "" && tbu.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request body received",
		})
	}

	q :=
		`
		query Balance($uid: string) {
			balance(func: uid($uid)) @normalize{
				name: name
				user {
					user_uid: uid
				}
			}
		}
		`

	resp, err := dgraph.NewTxn().QueryWithVars(context.Background(), q, map[string]string{"$uid": balanceId})

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	balance := struct {
		Result []struct {
			Name      string `json:"name"`
			UserUID string `json:"user_uid"`
		} `json:"balance"`
	}{}

	_ = json.Unmarshal(resp.Json, &balance)

	if len(balance.Result) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "The resource doesn't exist on this server",
		})
	}

	if balance.Result[0].UserUID != uid {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
		})
	}

	tbuJson, err := json.Marshal(tbu)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	mutation := &api.Mutation{
		CommitNow: true,
		SetJson:   tbuJson,
	}

	_, err = dgraph.NewTxn().Mutate(context.Background(), mutation)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Post updated successfully",
	})
}
