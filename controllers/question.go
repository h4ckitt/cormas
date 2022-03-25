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

func CreateQuestion(c *fiber.Ctx) error {
	question := new(models.Question)

	if err := c.BodyParser(question); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An Error Occurred While Processing That Request, Please Try Again Later",
		})
	}

	if question.Name == "" || question.Description == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Request	 Body Received",
		})
	}

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

	author := struct {
		UID string `json:"uid"`
	}{uid}

	question.Author = author
	question.Type = "Question"

	for index, tag := range *question.Tags {
		var (
			uid string
			err error
		)
		if uid, err = utils.GetTagUID(tag); err != nil {
			log.Println(err)
			if err.Error() == "tag doesn't exist yet" {
				newTag := models.HashTag{Type: "HashTag"}
				bytes, _ := json.Marshal(tag)
				_ = json.Unmarshal(bytes, &newTag)
				(*question.Tags)[index] = newTag
				continue
			}

			if err.Error() == "invalid tag" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"message": "Bad request body received",
				})
			}

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "An error occurred while processing that request",
			})
		}

		(*question.Tags)[index] = &models.HashTag{UID: uid}
	}

	questionJson, err := json.Marshal(question)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "An Error Occurred While Processing That Request,, Please Try Again Later",
		})
	}

	_, err = dgraph.NewTxn().Mutate(context.Background(), &api.Mutation{CommitNow: true, SetJson: questionJson})

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

func UpdateQuestion(c *fiber.Ctx) error {

	questionID := c.Params("id")

	if questionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request received",
		})
	}

	tbu := struct {
		UID         string `json:"uid"`
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`
		Privacy     int    `json:"privacy,omitempty"`
		UpdatedAt   string `json:"updated_at"`
	}{UID: questionID, UpdatedAt: time.Now().Format(time.RFC3339)}

	if err := c.BodyParser(&tbu); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	if tbu.Description == "" && tbu.Name == "" && tbu.Privacy == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request body received",
		})
	}

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

	q :=
		`
		query Questions($uid: string) {
			question(func: uid($uid)) @normalize {
				author {
					author_uid: uid
				}
			}
		}
		`

	resp, err := dgraph.NewTxn().QueryWithVars(context.Background(), q, map[string]string{"$uid": questionID})

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	authData := struct {
		Result []struct {
			AuthorUID string `json:"author_uid"`
		} `json:"question"`
	}{}

	json.Unmarshal(resp.Json, &authData)

	if len(authData.Result) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "The resource doesn't exist on this server",
		})
	}

	if authData.Result[0].AuthorUID != uid {
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
		"message": "Question updated successfully",
	})
}

func DeleteQuestion(c *fiber.Ctx) error {
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

	questionID := c.Params("id")

	if questionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid id received",
		})
	}

	q :=
		`
		query Questions($uid: string) {
			question(func: uid($uid)) @normalize {
				author {
					author_uid: uid
				}
			}
		}
		`

	resp, err := dgraph.NewTxn().QueryWithVars(context.Background(), q, map[string]string{"$uid": questionID})

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	authData := struct {
		Result []struct {
			AuthorUID string `json:"author_uid"`
		} `json:"question"`
	}{}

	json.Unmarshal(resp.Json, &authData)

	if len(authData.Result) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "The resource doesn't exist on this server",
		})
	}

	if authData.Result[0].AuthorUID != uid {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
		})
	}

	tbd := struct {
		UID string `json:"uid"`
	}{questionID}

	tbdJson, err := json.Marshal(tbd)

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

func ListQuestions(c *fiber.Ctx) error {
	q :=
		`
		questions(func: type(Question)) {
			uid
			name
			description
			author {
				name
				username
				email
				avatar
				cover
			}
			comments {
				uid
				description
				author {
					name
					username
					email
					avatar
					cover
				}
				reaction {
					name
				}
				updated_at
			}
			reactions {
				name
			}
			tags {
				name
			}
		}
		`

	questions := struct {
		Result []models.Question `json:"questions"`
	}{}

	resp, err := dgraph.NewTxn().Query(context.Background(), q)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An Error Occurred While Processing That Request",
		})
	}

	json.Unmarshal(resp.Json, &questions)

	return c.JSON(questions.Result)
}

func GetQuestion(c *fiber.Ctx) error {
	questionID := c.Params("id")

	if questionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request body received",
		})
	}

	q :=
		`
		query Questions($uid: string) {
			question(func: uid($uid)) {
				uid
				name
				description
				author {
					name
					username
					email
					avatar
					cover
				}
				comments {
					uid
					description
					author {
						name
						username
						email
						avatar
						cover
					}
					reaction {
						name
					}
					updated_at
				}
				reactions {
					name
				}
				tags {
					name
				}
			}
		}
		`

	resp, err := dgraph.NewTxn().QueryWithVars(context.Background(), q, map[string]string{"$uid": questionID})

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An Error Occurred While Processing That Request",
		})
	}

	question := struct {
		Result []models.Product `json:"question"`
	}{}

	json.Unmarshal(resp.Json, &question)

	if len(question.Result) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Requested resource was not found on this server",
		})
	}

	return c.JSON(question.Result[0])
}
