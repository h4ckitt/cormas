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

func CreateAsset(c *fiber.Ctx) error {
	uid, err := utils.GetJWTUser(c.Locals("user").(*jwt.Token))

	if err != nil {
		if err.Error() == "invalid JWT Token" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Forbidden",
			})
		}
	}

	asset := new(models.Asset)

	if err := c.BodyParser(asset); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request body was received",
		})
	}

	if asset.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request body was received",
		})
	}

	if asset.Document == "" && asset.Image == "" && asset.Video == "" && asset.Zip == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request body was received",
		})
	}

	postID := asset.EntityUID

	q :=
		`
		query Post($uid: string) {
			post(func: uid($uid)) @normalize{
				name: name
				author {
					author_uid: uid
				}
			}
		}
		`

	resp, err := dgraph.NewTxn().QueryWithVars(context.Background(), q, map[string]string{"$uid": postID})

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	post := struct {
		Result []struct {
			Name      string `json:"name"`
			AuthorUID string `json:"author_uid"`
		} `json:"post"`
	}{}

	err = json.Unmarshal(resp.Json, &post)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	if len(post.Result) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "The resource doesn't exist on this server",
		})
	}

	if post.Result[0].AuthorUID != uid {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
		})
	}

	asset.UID = "_:new"

	now := time.Now().Format(time.RFC3339)
	asset.CreatedAt, asset.UpdatedAt = now, now
	asset.EntityUID = "" //Set This To Empty Because It's Not Needed In The Database
	asset.Type = "Asset"

	commentJson, err := json.Marshal(asset)

	if err != nil {
		log.Println(err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	mutation := &api.Mutation{
		CommitNow: true,
		SetJson:   commentJson,
	}

	resp, err = dgraph.NewTxn().Mutate(context.Background(), mutation)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	postBody := struct {
		UID   string `json:"uid"`
		Asset []struct {
			UID string `json:"uid"`
		} `json:"assets"`
	}{UID: postID}

	x := struct {
		UID string `json:"uid"`
	}{resp.Uids["new"]}

	postBody.Asset = append(postBody.Asset, x)

	postBodyJson, err := json.Marshal(postBody)

	if err != nil {
		log.Println(err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	_, err = dgraph.NewTxn().Mutate(context.Background(), &api.Mutation{CommitNow: true, SetJson: postBodyJson})

	if err != nil {
		log.Println(err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Comment posted successfully",
	})
}

func DeleteAsset(c *fiber.Ctx) error {
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

	assetID := c.Params("id")

	q :=
		`
		query Asset($uid: string) {
			asset(func: uid($uid)) {
				author {
					author_uid: uid
				}
			}
		}
		`

	resp, err := dgraph.NewTxn().QueryWithVars(context.Background(), q, map[string]string{"$uid": assetID})

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	asset := struct {
		Result []struct {
			AuthorUID string `json:"author_uid"`
		} `json:"asset"`
	}{}

	err = json.Unmarshal(resp.Json, &asset)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	if len(asset.Result) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "The resource doesn't exist on this server",
		})
	}

	if asset.Result[0].AuthorUID != uid {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
		})
	}

	tbdJson, err := json.Marshal(struct {
		UID string `json:"uid"`
	}{uid})

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
