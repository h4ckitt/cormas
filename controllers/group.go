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

func CreateGroup(c *fiber.Ctx) error {
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

	group := new(models.Group)

	if err := c.BodyParser(group); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request body was received",
		})
	}

	if group.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request body was received",
		})
	}

	if group.Status < 1 || group.Status > 3 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request body was received",
		})
	}

	if group.Moderation < 1 || group.Moderation > 3 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request body was received",
		})
	}

	author := struct {
		UID string `json:"uid"`
	}{uid}

	now := time.Now().Format(time.RFC3339)

	group.Owner = author
	group.Editor = author
	group.Publisher = author
	group.CreatedAt, group.UpdatedAt = now, now
	group.Type = "Group"

	tx := struct {
		UID   string       `json:"uid"`
		Group models.Group `json:"posts"`
	}{uid, *group}

	groupJSon, err := json.Marshal(tx)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "An Error Occurred While Processing That Request,, Please Try Again Later",
		})
	}

	mutation := &api.Mutation{
		CommitNow: true,
		SetJson:   groupJSon,
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

func UpdateGroup(c *fiber.Ctx) error {
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

	groupID := c.Params("id")

	if groupID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request received",
		})
	}

	q := `
		query User($uid: string) {
			group(func: uid($uid)) @normalize{
				owner {
					owner_uid: uid
				}
			}
		}
		`

	resp, err := dgraph.NewTxn().QueryWithVars(context.Background(), q, map[string]string{"$uid": groupID})

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	group := struct {
		Result []struct {
			OwnerUID string `json:"owner_uid"`
		} `json:"group"`
	}{}

	err = json.Unmarshal(resp.Json, &group)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	if len(group.Result) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "The resource doesn't exist on this server",
		})
	}

	if group.Result[0].OwnerUID != uid {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
		})
	}

	tbu := struct {
		UID        string `json:"uid"`
		Name       string `json:"name"`
		About      string `json:"about,omitempty"`
		Status     int    `json:"status"`
		Moderation int    `json:"moderation"`
		ParentUID  string `json:"child,omitempty"`
		UpdatedAt  string `json:"updated_at"`
	}{UID: groupID, UpdatedAt: time.Now().Format(time.RFC3339)}

	if err := c.BodyParser(&tbu); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	if tbu.Name == "" || tbu.Status < 1 || tbu.Status > 3 || tbu.Moderation < 1 || tbu.Moderation > 3 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request body received",
		})
	}

	q = `
		query User($uid: string) {
			group(func: uid($uid)) @normalize{
				owner {
					owner_uid: uid
				}
			}
		}
	`
	resp, err = dgraph.NewTxn().QueryWithVars(context.Background(), q, map[string]string{"$uid": tbu.ParentUID})
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}
	groupTwo := struct {
		Result []struct {
			OwnerUID string `json:"owner_uid"`
		} `json:"group"`
	}{}

	err = json.Unmarshal(resp.Json, &groupTwo)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	if len(groupTwo.Result) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "The resource doesn't exist on this server",
		})
	}

	if groupTwo.Result[0].OwnerUID != uid {
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
		"message": "Comment updated successfully",
	})
}

func DeleteGroup(c *fiber.Ctx) error {
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

	groupId := c.Params("id")

	q :=
		`
		query Group($uid: string) {
			group(func: uid($uid)) @normalize{
				owner {
					owner_uid: uid
				}
			}
		}
		`

	resp, err := dgraph.NewTxn().QueryWithVars(context.Background(), q, map[string]string{"$uid": groupId})

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	group := struct {
		Result []struct {
			OwnerUID string `json:"owner_uid"`
		} `json:"group"`
	}{}

	err = json.Unmarshal(resp.Json, &group)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	if len(group.Result) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "The resource doesn't exist on this server",
		})
	}

	if group.Result[0].OwnerUID != uid {
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

func SuspendGroup(c *fiber.Ctx) error {
	return nil
}

func ListGroups(c *fiber.Ctx) error {
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
		query User($uid: string) {
			groups(func: uid($uid)) {
				groups {
					name
					icon
					about
					owner {
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

	resp, err := dgraph.NewTxn().QueryWithVars(context.Background(), q, map[string]string{"$uid": uid})

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	groups := struct {
		Result []models.Comment `json:"groups"`
	}{}

	err = json.Unmarshal(resp.Json, &groups)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	return c.Status(fiber.StatusOK).JSON(groups)
}
