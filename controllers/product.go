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
)

func CreateProduct(c *fiber.Ctx) error {
	_, err := utils.GetJWTUser(c.Locals("user").(*jwt.Token))

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

	product := new(models.Product)

	if err := c.BodyParser(&product); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Baf request body received",
		})
	}

	if product.Name == "" || product.RegularPrice <= 0 || product.SellingPrice <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request body received",
		})
	}

	//Insert business as owner of product once we figure out how to get author's current location

	product.Type = "Product"

	productJson, err := json.Marshal(product)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	_, err = dgraph.NewTxn().Mutate(context.Background(), &api.Mutation{CommitNow: true, SetJson: productJson})

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(productJson)
}

func UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request body received",
		})
	}

	tbu := struct {
		UID                   string  `json:"uid"`
		Name                  string  `json:"name,omitempty"`
		Description           string  `json:"description,omitempty"`
		RegularPrice          float64 `json:"regular_price,omitempty"`
		SellingPrice          float64 `json:"selling_price,omitempty"`
		ProductType           int     `json:"type,omitempty"`
		Supported             int     `json:"supported,omitempty"`
		Excerpt               string  `json:"excerpt,omitempty"`
		TechnicalInformation  string  `json:"technical_information,omitempty"`
		AdditionalInformation string  `json:"additional_information,omitempty"`
		ProductGuides         string  `json:"product_guides,omitempty"`
	}{}

	if err := c.BodyParser(&tbu); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "bad request body received",
		})
	}

	if tbu.Name == "" && tbu.Description == "" && tbu.RegularPrice <= 0 && tbu.SellingPrice <= 0 && tbu.ProductType == 0 &&
		tbu.Supported == 0 && tbu.Excerpt == "" && tbu.TechnicalInformation == "" && tbu.AdditionalInformation == "" &&
		tbu.ProductGuides == "" {
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
		query Product($uid: string) {
			product(func: uid($uid)) @normalize {
				owner {
					owner_uid: uid
				}
			}
		}
		`

	resp, err := dgraph.NewTxn().QueryWithVars(context.Background(), q, map[string]string{"$uid": id})

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	authData := struct {
		Result []struct {
			OwnerUid string `json:"uid"`
		} `json:"product"`
	}{}

	json.Unmarshal(resp.Json, &authData)

	if len(authData.Result) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "The resource doesn't exist on this server",
		})
	}

	if authData.Result[0].OwnerUid != uid {
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
		"message": "Post updated successfully",
	})
}

func DeleteProduct(c *fiber.Ctx) error {
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

	productID := c.Params("id")

	q :=
		`
		query Product($uid: string) {
			product(func: uid($uid)) @normalize {
				owner {
					owner_uid: uid
				}
			}
		}
		`

	resp, err := dgraph.NewTxn().QueryWithVars(context.Background(), q, map[string]string{"$uid": productID})

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An error occurred while processing that request",
		})
	}

	authData := struct {
		Result []struct {
			OwnerUid string `json:"uid"`
		} `json:"product"`
	}{}

	json.Unmarshal(resp.Json, &authData)

	if len(authData.Result) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "The resource doesn't exist on this server",
		})
	}

	if authData.Result[0].OwnerUid != uid {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
		})
	}

	tbd := struct {
		UID string `json:"uid"`
	}{productID}

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

func ListProducts(c *fiber.Ctx) error {
	q :=
		`
		products(func: type(Product)) {
			uid
			name
			description
			regular_price
			selling_price
			currency {
				name
				icon
				value
			}
			category {
				name
				icon
			}
			reviews {
				name
				description
				author {
					name
					username
					cover
					avatar
					email
				}
				rating
			}
			type
			supported
			downloadable {
				name
				image
				video
				document
				zip
			}
			thumbnail {
				image
			}
			gallery {
				name
				image
				video
				document
				zip
			}
			excerpt
			description
			technical_information
			additional_information
			product_information
			product_guides
			owner {
				name
				email
				username
				cover
				avatar
			}
		}
		`

	products := struct {
		Result []models.Product `json:"products"`
	}{}

	resp, err := dgraph.NewTxn().Query(context.Background(), q)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An Error Occurred While Processing That Request",
		})
	}

	json.Unmarshal(resp.Json, &products)

	return c.JSON(products.Result)
}

func GetProduct(c *fiber.Ctx) error {
	productID := c.Params("id")

	q :=
		`
		query Products($uid: string) {
			product(func: uid($uid)) {
				uid
				name
				description
				regular_price
				selling_price
				currency {
					name
					icon
					value
				}
				category {
					name
					icon
				}
				reviews {
					name
					description
					author {
						name
						username
						cover
						avatar
						email
					}
					rating
				}
				type
				supported
				downloadable {
					name
					image
					video
					document
					zip
				}
				thumbnail {
					image
				}
				gallery {
					name
					image
					video
					document
					zip
				}
				excerpt
				description
				technical_information
				additional_information
				product_information
				product_guides
				owner {
					name
					email
					username
					cover
					avatar
				}
			}
		}
		`

	resp, err := dgraph.NewTxn().QueryWithVars(context.Background(), q, map[string]string{"$uid": productID})

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "An Error Occurred While Processing That Request",
		})
	}

	product := struct {
		Result []models.Product `json:"product"`
	}{}

	json.Unmarshal(resp.Json, &product)

	if len(product.Result) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Requested resource was not found on this server",
		})
	}

	return c.JSON(product.Result[0])

}
