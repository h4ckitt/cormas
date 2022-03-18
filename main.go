package main

import (
	"flag"
	"log"
	"rest-api/db"
	"rest-api/routes"
)

func main() {

	createSchema := flag.Bool("create", false, "Create Database Schema")
	dropTable := flag.Bool("drop", false, "Truncate Database")

	flag.Parse()

	if *dropTable {
		db.DropAll()
	}

	if *createSchema {
		db.CreateSchema()

	}
	
	if err := routes.InitRouter().Listen(":8081"); err != nil {
		log.Fatal(err)
	}
}
