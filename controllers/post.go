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

func CreatePost(c *fiber.Ctx) error {
	uid := utils.GetJWTUser(c.Locals("user").(*jwt.Token))

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
	postID := c.Params("id")

	q :=
		`
		query Post($uid: string){
			post(func: uid($uid)){
				name
				description
				reactions {
					name
				}
				assets {
					image
					video
					document
				}
				author {
					name
					username
					email
					avatar
					cover
				}
				privacy
				address {
					name
					address1
				}
				comments {
					description
					created_at
					updated_at
					author {
						name
						username
						email
						avatar
						cover
					}
					reply {
						description
						created_at
						updated_at
						author {
							name
							username
							email
							avatar
							cover
						}
					}
				}
				tags {
					name
				}
			}
		}
		`

	resp, err := dgraph.NewTxn().QueryWithVars(context.Background(), q, map[string]string{"$uid": postID})

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An Error Occurred While Processing That Request",
		})
	}

	post := struct {
		Result []models.Post `json:"post"`
	}{}

	json.Unmarshal(resp.Json, &post)

	if len(post.Result) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Requested resource was not found on this server",
		})
	}

	return c.JSON(post)
}

func DeletePost(c *fiber.Ctx) error {
	uid := utils.GetJWTUser(c.Locals("user").(*jwt.Token))

	postID := c.Params("id")

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

	tbd := struct {
		UID string `json:"uid"`
	}{postID}

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

func UpdatePost(c *fiber.Ctx) error {

	uid := utils.GetJWTUser(c.Locals("user").(*jwt.Token))

	postID := c.Params("id")

	if postID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request received",
		})
	}

	tbu := struct {
		UID         string `json:"uid"`
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`
		UpdatedAt   string `json:"updated_at"`
	}{UID: postID, UpdatedAt: time.Now().Format(time.RFC3339)}

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
