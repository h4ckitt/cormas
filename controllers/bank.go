package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"rest-api/models"
	"rest-api/utils"
	"time"
)

func CreateBank(c *fiber.Ctx) error {
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

	balance := new(models.Bank)

	if err := c.BodyParser(&balance); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An Error Occurred While Processing That Request, Please Try Again Later",
		})
	}

	if balance.Name == "" || balance.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Request Body Received",
		})
	}

	/*author := struct {
		UID string `json:"uid"`
	}{uid}*/
	now := time.Now().Format(time.RFC3339)

	balance.CreatedAt = now
	balance.UpdatedAt = now
	balance.Type = "Bank"

	userBank := struct {
		UID  string      `json:"uid"`
		Bank models.Bank `json:"bank"`
	}{uid, *balance}

	userJson, err := json.Marshal(userBank)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	fmt.Println(string(userJson))

	_, err = dgraph.NewTxn().Mutate(context.Background(), &api.Mutation{CommitNow: true, SetJson: userJson})

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Bank Created Successfully",
	})
}

func GetBank(c *fiber.Ctx) error {
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
		query Bank($uid: string) {
			bank(func: uid($uid)){
				name
				amount
				currency {
					name
					value
					icon
					status
				}
				user {
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
		"$uid": balanceID,
	}

	resp, err := dgraph.NewTxn().QueryWithVars(context.Background(), q, variables)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An Error Occurred While Processing That Request",
		})
	}

	balance := struct {
		Result []models.Bank `json:"bank"`
	}{}

	_ = json.Unmarshal(resp.Json, &balance)

	if len(balance.Result) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Requested resource was not found on this server",
		})
	}

	if balance.Result[0].UID != uid {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
		})
	}

	return c.JSON(balance.Result[0])
}

func DeleteBank(c *fiber.Ctx) error {
	uid, _ := utils.GetJWTUser(c.Locals("user").(*jwt.Token))

	balanceId := c.Params("id")

	q :=
		`
		query Bank($uid: string) {
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
			Name    string `json:"name"`
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

func UpdateBank(c *fiber.Ctx) error {
	uid, _ := utils.GetJWTUser(c.Locals("user").(*jwt.Token))

	balanceId := c.Params("id")

	if balanceId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request received",
		})
	}

	tbu := struct {
		UID       string  `json:"uid"`
		Name      string  `json:"name,omitempty"`
		Amount    float64 `json:"amount,omitempty"`
		UpdatedAt string  `json:"updated_at"`
	}{UID: balanceId, UpdatedAt: time.Now().Format(time.RFC3339)}

	if err := c.BodyParser(&tbu); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	if tbu.Amount <= 0 && tbu.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request body received",
		})
	}

	q :=
		`
		query Bank($uid: string) {
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
			Name    string `json:"name"`
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
		"message": "Bank updated successfully",
	})
}
