package main

import (
	"fmt"
	"log"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"google.golang.org/grpc"
)

func newClient() *dgo.Dgraph {
	d, err := grpc.Dial("localhost:9080", grpc.WithInsecure())

	if err != nil {
		log.Fatal(err)
	}

	return dgo.NewDgraphClient(api.NewDgraphClient(d))
}

func main() {
	fmt.Println("vim-go")
}
