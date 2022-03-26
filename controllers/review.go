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

func CreateReview(c *fiber.Ctx) error {
	var (
		q        string
		dataType interface{}
		Type     string
	)
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

	review := new(models.Review)

	if err := c.BodyParser(&review); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request body received",
		})
	}

	if review.Name == "" || review.Description == "" || (review.Business == nil && review.Product == nil) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request body received",
		})
	}

	review.Type = "Review"
	review.Author = struct {
		UID string `json:"uid"`
	}{uid}
	now := time.Now().Format(time.RFC3339)
	review.CreatedAt, review.UpdatedAt = now, now

	tx := struct {
		UID    string        `json:"uid"`
		Review models.Review `json:"reviews"`
	}{Review: *review}

	if review.Business != nil {
		dataType = *(review.Business)
		Type = "User"

	} else if review.Product != nil {
		dataType = *(review.Product)
		Type = "Product"
	}

	if uid, ok := dataType.(string); ok {
		tx.UID = uid

		switch Type {
		case "User":
			*(review.Business) = struct {
				UID string `json:"uid"`
			}{uid}

		case "Product":
			*(review.Product) = struct {
				UID string `json:"uid"`
			}{uid}
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request body received",
		})
	}

	q =
		`
		query DB($uid: string) {
			entity(func: uid($uid)) {
				dgraph.type
			}
		}
		`

	resp, err := dgraph.NewTxn().QueryWithVars(context.Background(), q, map[string]string{"$uid": tx.UID})

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	res := struct {
		Result []struct {
			Type []string `json:"dgraph.type"`
		} `json:"entity"`
	}{}

	_ = json.Unmarshal(resp.Json, &res)

	if len(res.Result) == 0 || res.Result[0].Type[0] != Type {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "The resource doesn't exist on this server",
		})
	}

	reviewJson, err := json.Marshal(tx)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	_, err = dgraph.NewTxn().Mutate(context.Background(), &api.Mutation{CommitNow: true, SetJson: reviewJson})

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Review Created Successfully",
	})
}

func UpdateReview(c *fiber.Ctx) error {
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

	reviewID := c.Params("id")

	if reviewID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad response body received",
		})
	}

	tbu := struct {
		UID         string `json:"uid"`
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`
		Rating      int    `json:"rating,omitempty"`
		UpdatedAt   string `json:"updated_at"`
	}{UID: reviewID, Rating: -1, UpdatedAt: time.Now().Format(time.RFC3339)}

	if err := c.BodyParser(&tbu); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request body received",
		})
	}

	if tbu.Description == "" && tbu.Name == "" && tbu.Rating < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request body received",
		})
	}

	q :=
		`
		query Reviews($uid: string) {
			review(func: uid($uid)) @normalize{
				author {
					author_uid: uid
				}
			}
		}
		`

	resp, err := dgraph.NewTxn().QueryWithVars(context.Background(), q, map[string]string{"$uid": reviewID})

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	res := struct {
		Result []struct {
			AuthorUID string `json:"author_uid"`
		} `json:"review"`
	}{}

	err = json.Unmarshal(resp.Json, &res)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	if len(res.Result) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "The resource doesn't exist on this server",
		})
	}

	if res.Result[0].AuthorUID != uid {
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

	_, err = dgraph.NewTxn().Mutate(context.Background(), &api.Mutation{CommitNow: true, SetJson: tbuJson})

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Review updated successfully",
	})
}

func DeleteReview(c *fiber.Ctx) error {
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

	reviewID := c.Params("id")

	q :=
		`
		query Reviews($uid: string) {
			review(func: uid($uid)) @normalize{
				author {
					author_uid: uid
				}
			}
		}
		`

	resp, err := dgraph.NewTxn().QueryWithVars(context.Background(), q, map[string]string{"$uid": reviewID})

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	res := struct {
		Result []struct {
			AuthorUID string `json:"author_uid"`
		} `json:"review"`
	}{}

	err = json.Unmarshal(resp.Json, &res)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	if len(res.Result) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "The resource doesn't exist on this server",
		})
	}

	if res.Result[0].AuthorUID != uid {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
		})
	}

	tbdJson, err := json.Marshal(struct {
		UID string `json:"uid"`
	}{reviewID})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	_, err = dgraph.NewTxn().Mutate(context.Background(), &api.Mutation{CommitNow: true, DeleteJson: tbdJson})

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
