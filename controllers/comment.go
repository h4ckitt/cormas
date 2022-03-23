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

	if comment.Description == "" || comment.EntityUID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request body was received",
		})
	}

	postID := comment.EntityUID

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

	json.Unmarshal(resp.Json, &post)

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

	comment.UID = "_:new"
	comment.Author = struct {
		UID string `json:"uid"`
	}{uid}
	now := time.Now().Format(time.RFC3339)
	comment.CreatedAt, comment.UpdatedAt = now, now
	comment.EntityUID = "" //Set This To Empty Because It's Not Needed In The Database
	comment.Type = "Comment"

	commentJson, err := json.Marshal(comment)

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
		UID     string `json:"uid"`
		Comment []struct {
			UID string `json:"uid"`
		} `json:"comments"`
	}{UID: postID}

	x := struct {
		UID string `json:"uid"`
	}{resp.Uids["new"]}

	postBody.Comment = append(postBody.Comment, x)

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

	json.Unmarshal(resp.Json, &comment)

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

	json.Unmarshal(resp.Json, &comment)

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

	json.Unmarshal(resp.Json, &comments)

	return c.Status(fiber.StatusOK).JSON(comments)
}
