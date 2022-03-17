package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"rest-api/db"
	"rest-api/models"
	"rest-api/utils"
	"time"
)

var dgraph = db.GetDB()

func SignUpHandler(c *fiber.Ctx) error {

	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An Error Occurred, Please Try Again Later",
		})
	}

	if user.Name == "" || user.Password == "" || user.Email == "" || user.Username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "bad request body received",
		})
	}

	mutation := &api.Mutation{CommitNow: true}

	user.Type = "User"
	now := time.Now().Format(time.RFC3339)
	user.CreatedAt = now
	user.UpdatedAt = now
	userJson, err := json.Marshal(user)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An Error Occurred, Please Try Again Later",
		})
	}

	mutation.SetJson = userJson

	_, err = dgraph.NewTxn().Mutate(context.Background(), mutation)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "An Error Occurred, Please Try Again Later",
		})
	}

	//password, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 12)

	user.Password = ""

	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "An Error Occurred, Please Try Again",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

func Login(c *fiber.Ctx) error {

	loginData := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	authData := struct {
		Result []struct {
			Uid   string `json:"uid"`
			Valid bool   `json:"validPass"`
		} `json:"auth"`
	}{}

	if err := c.BodyParser(&loginData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	if loginData.Email == "" || loginData.Password == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "Invalid Message Body Received",
		})
	}

	fmt.Println(loginData.Password)
	variables := map[string]string{"$email": loginData.Email, "$pass": loginData.Password}

	q :=
		`
		query User($email: string, $pass: string){
			auth(func: type(User)) @filter(eq(email, $email)) {
				uid
				validPass: checkpwd(password,$pass)
			}
		}
	`

	resp, err := dgraph.NewTxn().QueryWithVars(context.Background(), q, variables)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "An Error Occurred While Processing That Request, Please Try Again",
		})
	}

	err = json.Unmarshal(resp.Json, &authData)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An Error Occurred While Processing The Request, Please Try Again",
		})
	}

	if !authData.Result[0].Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Email Or Password Invalid",
		})
	}

	token, err := utils.GenerateJWTCookie(authData.Result[0].Uid)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An Error Occurred, Please Try Again Later",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "authorization_token",
		Value:    token,
		Expires:  time.Now().Add(1 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User Logged In Successfully",
	})
}

func Update(c *fiber.Ctx) error {
	uid := utils.GetJWTUser(c.Locals("user").(*jwt.Token))

	updatedUser := struct {
		UID       string `json:"uid"`
		Cover     string `json:"cover,omitempty"`
		Avatar    string `json:"avatar,omitempty"`
		Username  string `json:"username,omitempty"`
		UpdatedAt string `json:"updated_at"`
	}{}

	if err := c.BodyParser(&updatedUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request Body Request Received",
		})
	}

	if updatedUser.Cover == "" && updatedUser.Avatar == "" && updatedUser.Username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request Body Request Received",
		})
	}

	updatedUser.UID = uid
	updatedUser.UpdatedAt = time.Now().Format(time.RFC3339)

	updatedJsonBody, err := json.Marshal(updatedUser)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An Error Occurred While Processing That Request, Please Try Again Later",
		})
	}

	mutation := &api.Mutation{
		CommitNow: true,
		SetJson:   updatedJsonBody,
	}

	_, err = dgraph.NewTxn().Mutate(context.Background(), mutation)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An Error Occurred While Processing This Request, Please Try Again",
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "User Updated Successfully",
	})
}
