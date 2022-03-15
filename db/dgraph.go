package db

import (
	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"google.golang.org/grpc"
	"log"
)

var db *dgo.Dgraph

func init() {
	conn, err := grpc.Dial("localhost:9080", grpc.WithInsecure())

	if err != nil {
		log.Fatal(err)
	}

	db = dgo.NewDgraphClient(api.NewDgraphClient(conn))

	op := &api.Operation{}

	op.Schema = `
		name: string @index(exact) .
		email: string @index(hash) .
		username: string @index(hash) .
		password: string .
		verified: bool .
		LastIP: string .
		Premium: int .
		isBusiness: bool .
		moderation: int .
		privacy: int .
		Follows: [uid] @count .
		
		type User {
			name: string
			email: string
			username: string
			password: string
			address: [
`
}

func GetDB() *dgo.Dgraph {
	return db
}
