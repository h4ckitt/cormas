package main

import (
	"log"
	"rest-api/routes"
)

func main() {
	if err := routes.InitRouter().Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
