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

func CreateComment(c *fiber.Ctx) error {
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

	comment := new(models.Comment)

	if err := c.BodyParser(comment); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request body was received",
		})
	}

	if comment.Description == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request body was received",
		})
	}

	comment.Author = struct {
		UID string `json:"uid"`
	}{uid}
	now := time.Now().Format(time.RFC3339)
	comment.CreatedAt, comment.UpdatedAt = now, now
	comment.Type = "Comment"

	tx := struct {
		UID     string         `json:"uid"`
		Comment models.Comment `json:"comments"`
	}{Comment: *comment}

	if comment.Post != nil {
		dataType = *(comment.Post)
		Type = "Post"
	} else if comment.Question != nil {
		dataType = *(comment.Question)
		Type = "Question"
	}

	if uid, ok := dataType.(string); ok {
		tx.UID = uid

		switch Type {
		case "Post":
			*(tx.Comment.Post) = struct {
				UID string `json:"uid"`
			}{uid}

		case "Question":
			*(tx.Comment.Question) = struct {
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
				name
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

	post := struct {
		Result []struct {
			Name string   `json:"name"`
			Type []string `json:"dgraph.type"`
		} `json:"entity"`
	}{}

	err = json.Unmarshal(resp.Json, &post)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	if len(post.Result) == 0 || post.Result[0].Type[0] != Type {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "The resource doesn't exist on this server",
		})
	}

	commentJson, err := json.Marshal(tx)

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

	_, err = dgraph.NewTxn().Mutate(context.Background(), mutation)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

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

func UpdateComment(c *fiber.Ctx) error {

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

	commentID := c.Params("id")

	if commentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad response body received",
		})
	}

	q :=
		`
		query Post($uid: string) {
			comment(func: uid($uid)) @normalize{
				author {
					author_uid: uid
				}
			}
		}
		`

	resp, err := dgraph.NewTxn().QueryWithVars(context.Background(), q, map[string]string{"$uid": commentID})

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	comment := struct {
		Result []struct {
			AuthorUID string `json:"author_uid"`
		} `json:"comment"`
	}{}

	err = json.Unmarshal(resp.Json, &comment)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	if len(comment.Result) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "The resource doesn't exist on this server",
		})
	}

	if comment.Result[0].AuthorUID != uid {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
		})
	}

	tbu := struct {
		UID         string `json:"uid"`
		Description string `json:"description"`
		UpdatedAt   string `json:"updated_at"`
	}{UID: commentID, UpdatedAt: time.Now().Format(time.RFC3339)}

	if err := c.BodyParser(&tbu); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request body received",
		})
	}

	if tbu.Description == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request body received",
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
		"message": "Comment updated successfully",
	})
}

func DeleteComment(c *fiber.Ctx) error {
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

	commentID := c.Params("id")

	q :=
		`
		query Comment($uid: string) {
			comment(func: uid($uid)) @normalize{
				author {
					author_uid: uid
				}
			}
		}
		`

	resp, err := dgraph.NewTxn().QueryWithVars(context.Background(), q, map[string]string{"$uid": commentID})

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	comment := struct {
		Result []struct {
			AuthorUID string `json:"author_uid"`
		} `json:"comment"`
	}{}

	err = json.Unmarshal(resp.Json, &comment)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	if len(comment.Result) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "The resource doesn't exist on this server",
		})
	}

	if comment.Result[0].AuthorUID != uid {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
		})
	}

	tbdJson, err := json.Marshal(struct {
		UID string `json:"uid"`
	}{commentID})

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

func ListComments(c *fiber.Ctx) error {
	postID := c.Params("id")

	q :=
		`
		query Post($uid: string) {
			comments(func: uid($uid)) {
				comments {
					description
					author {
						name
						username
						email
						avatar
						cover
					}
					
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

	comments := struct {
		Result []models.Comment `json:"comments"`
	}{}

	err = json.Unmarshal(resp.Json, &comments)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	return c.Status(fiber.StatusOK).JSON(comments)
}
