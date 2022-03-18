package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"github.com/gofiber/fiber/v2"
	"log"
	"rest-api/models"
	"time"
)

func CreateCurrency(c *fiber.Ctx) error {

	currency := new(models.Currency)

	if err := c.BodyParser(currency); err != nil {
		fmt.Println("An Error Occurred: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	if currency.Name == "" || currency.Value == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request Body",
		})
	}

	mu := &api.Mutation{CommitNow: true}

	currency.Type = "Currency"
	now := time.Now().Format(time.RFC3339)
	currency.CreatedAt, currency.CreatedAt = now, now

	currencyBody, err := json.Marshal(currency)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "An Error Occurred. Please Try Again",
		})
	}

	mu.SetJson = currencyBody

	_, err = dgraph.NewTxn().Mutate(context.Background(), mu)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "An Error Occurred",
		})
	}

	return c.Status(fiber.StatusOK).JSON(currency)
}

func GetAllCurrencies(c *fiber.Ctx) error {
	q := `
		query Currency() {
			currency(func: type(Currency)){
				uid
				name
				value
				status
			  }
		}`

	resp, err := dgraph.NewTxn().Query(context.Background(), q)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred. Please try again",
		})
	}

	currencyData := struct {
		Result  []struct{
			Uid string `json:"uid"`
			Name string `json:"name"`
			Value string `json:"value"`
			Status int	`json:"status"`
		} `json:"currency"`
	}{}

	err = json.Unmarshal(resp.Json, &currencyData)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred. Please try again",
		})
	}

	return c.Status(fiber.StatusOK).JSON(currencyData)
}

func GetOneCurrency(c *fiber.Ctx) error {
	uidToGet := c.Params("uid")
	if uidToGet == "" {
		return 	c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "ID not found",
		})
	}

	variable := map[string]string{"$uid": uidToGet}

	q := `
		query Currency($uid: string) {
			currency(func: uid($uid)) {
				uid
				name
				value
				status
			  }
		}`

	resp, err := dgraph.NewTxn().QueryWithVars(context.Background(), q, variable)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "An error occurred. Please try again",
		})
	}

	currencyData := struct {
		Result  []struct{
			Uid string `json:"uid"`
			Name string `json:"name"`
			Value string `json:"value"`
			Status int	`json:"status"`
		} `json:"currency"`
	}{}

	err = json.Unmarshal(resp.Json, &currencyData)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred. Please try again",
		})
	}

	return c.Status(fiber.StatusOK).JSON(currencyData)
}

func UpdateCurrency(c *fiber.Ctx) error {
	return nil
}

func DeleteCurrency(c *fiber.Ctx) error {
	return nil
}
