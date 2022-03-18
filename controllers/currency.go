package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/gofiber/fiber/v2"
	"log"
	"rest-api/db"
	"rest-api/models"
)

func CreateCurrency(c *fiber.Ctx) error {
	dgraph := db.GetDB()

	txn := dgraph.NewTxn()

	defer func(txn *dgo.Txn, ctx context.Context) {
		err := txn.Discard(ctx)
		if err != nil {
			log.Println(err)
		}
	}(txn, context.TODO())

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

	currencyBody, err := json.Marshal(currency)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "An Error Occurred. Please Try Again",
		})
	}

	mu := &api.Mutation{
		SetJson: currencyBody,
	}

	mu.CommitNow = true

	_, err = txn.Mutate(context.TODO(), mu)

	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "An Error Occurred",
		})
	}

	return c.Status(fiber.StatusOK).JSON(currency)
}

func GetAllCurrencies(c *fiber.Ctx) error {
	return nil
}

func UpdateCurrency(c *fiber.Ctx) error {
	return nil
}

func DeleteCurrency(c *fiber.Ctx) error {
	return nil
}