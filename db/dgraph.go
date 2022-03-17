package db

import (
	"context"
	"fmt"

	//	"github.com/fenos/dqlx"
	//"github.com/dgraph-io/dgo"
	//"github.com/dgraph-io/dgo/protos/api"
	"github.com/fenos/dqlx"
	//"google.golang.org/grpc"
	"log"
)

//var db *dgo.Dgraph
var db dqlx.DB

func CreateSchema(drop bool) {
	schema := db.Schema()

	if drop {
		schema.Alter(context.Background(), dqlx.WithDropAllSchema(true))
	}

	name := schema.Predicate("name", dqlx.ScalarString)
	moderation := schema.Predicate("moderation", dqlx.ScalarBool)
	userAgent := schema.Predicate("user_agent", dqlx.ScalarString)
	privacy := schema.Predicate("privacy", dqlx.ScalarInt)
	amount := schema.Predicate("amount", dqlx.ScalarFloat)
	createdAt := schema.Predicate("created_at", dqlx.ScalarDateTime)
	updatedAt := schema.Predicate("updated_at", dqlx.ScalarDateTime)
	status := schema.Predicate("status", dqlx.ScalarInt)
	icon := schema.Predicate("icon", dqlx.ScalarString)
	description := schema.Predicate("description", dqlx.ScalarString)

	schema.Type("User", func(user *dqlx.TypeBuilder) {
		user.Predicate(name)
		user.String("username").IndexExact()
		user.Password("password")
		user.String("avatar")
		user.String("cover")
		user.Predicate(moderation)
		user.Bool("is_business")
		user.Bool("verified")
		user.Int("premium")
		//user.Float("amount")
		user.Predicate(amount)
		user.String("last_ip")
		user.String("user_agent")
		user.Float("balance")
		user.Type("address", "Address").List()
		user.Type("currency", "Currency")
		user.Type("bank", "Balance")
		user.Type("orders", "Order").List()
		user.Type("invoices", "Invoice").List()
		user.Type("posts", "Post").List()
		user.Type("reviews", "Review").List()
		user.Type("owners", "User").List()
		user.Type("editor", "User")
		user.Type("publisher", "User")
		user.Type("category", "Category")
		user.Type("sales", "Order").List()
		user.Type("sale_invoices", "Invoice").List()
		user.Predicate(privacy)
		user.Predicate(createdAt)
		user.Predicate(updatedAt)

	})

	schema.Type("Post", func(post *dqlx.TypeBuilder) {
		post.Predicate(name)
		post.Predicate(description)
		post.Predicate(privacy)
		post.Predicate(moderation)
		post.Predicate(amount)
		post.Type("author", "User").Reverse()
		post.Type("business", "User")
		post.Type("orders", "Order").List()
		post.Type("address", "Address")
		post.Type("currency", "Currency")
		post.Type("comments", "Comment").List()
		post.Type("reactions", "Reaction").List()
		post.Type("tags", "HashTag").List()
		post.Type("assets", "Asset").List()

	})

	schema.Type("Order", func(order *dqlx.TypeBuilder) {
		order.Predicate(status)
		order.Predicate(amount)
		order.Predicate(moderation)
		order.Predicate(userAgent)
		order.Predicate(createdAt)
		order.Predicate(updatedAt)
		order.String("transaction_id")
		order.Type("business", "Business")
		order.Type("sender", "Sender")
		order.Type("receiver", "Receiver")
		order.Type("currency", "Currency")
		order.Type("posts", "Post").List()
		order.Type("products", "Product").List()
	})

	schema.Type("Address", func(address *dqlx.TypeBuilder) {
		address.Predicate(name)
		address.Predicate(status)
		address.Predicate(moderation)
		address.String("address_1")
		address.String("address_2")
		address.String("city")
		address.String("state")
		address.String("country")
	})

	schema.Type("Asset", func(asset *dqlx.TypeBuilder) {
		asset.Predicate(name)
		asset.Predicate(createdAt)
		asset.Predicate(updatedAt)
		asset.Predicate(moderation)
		asset.String("image")
		asset.String("video")
		asset.String("document")
		asset.String("zip")
		asset.Type("post", "Post")
		asset.Type("comment", "Comment")
		asset.Type("review", "Review")
	})

	schema.Type("Currency", func(currency *dqlx.TypeBuilder) {
		currency.Predicate(name)
		currency.Predicate(icon)
		currency.Predicate(status)
		currency.Predicate(createdAt)
		currency.Predicate(updatedAt)
		currency.String("value")
	})

	schema.Type("Balance", func(balance *dqlx.TypeBuilder) {
		balance.Predicate(amount)
		balance.Predicate(status)
		balance.Predicate(moderation)
		balance.Predicate(createdAt)
		balance.Predicate(updatedAt)
		balance.Type("currency", "Currency")
	})

	schema.Type("Category", func(category *dqlx.TypeBuilder) {
		category.Predicate(name)
		category.Predicate(icon)
		category.Predicate(moderation)
		category.Predicate(status)
		category.Type("child", "Category").List()
	})

	schema.Type("Comment", func(comment *dqlx.TypeBuilder) {
		comment.Predicate(description)
		comment.Predicate(moderation)
		comment.Predicate(createdAt)
		comment.Predicate(updatedAt)
		comment.Type("author", "User")
		comment.Type("post", "Post").Reverse()
		comment.Type("reactions", "Reaction")
		comment.Type("replies", "Comment").List()
		comment.Type("address", "Address")
	})

	schema.Type("Invoice", func(invoice *dqlx.TypeBuilder) {
		invoice.Predicate(status)
		invoice.Predicate(createdAt)
		invoice.Predicate(updatedAt)
		invoice.Type("order", "Order")
		invoice.Type("user", "buyer")
	})

	schema.Type("Product", func(product *dqlx.TypeBuilder) {
		product.Predicate(name)
		product.Predicate(description)
		product.Predicate(moderation)
		product.Predicate(status)
		product.Int("supported")
		product.Int("type")
		product.String("excerpt")
		product.String("technical_information")
		product.String("additional_information")
		product.String("product_information")
		product.String("product_guides")
		product.Float("regular_price")
		product.Float("selling_price")
		product.Type("address", "Address")
		product.Type("owner", "User")
		product.Type("currency", "Currency")
		product.Type("category", "Category")
		product.Type("reviews", "Review").List()
		product.Type("thumbnail", "Asset")
		product.Type("downloadable", "Asset").List()
		product.Type("gallery", "Asset").List()
	})

	schema.Type("Question", func(question *dqlx.TypeBuilder) {
		question.Predicate(name)
		question.Predicate(description)
		question.Predicate(status)
		question.Predicate(moderation)
		question.Predicate(privacy)
		question.Predicate(createdAt)
		question.Predicate(updatedAt)
		question.Type("tags", "HashTag").List()
		question.Type("comments", "Comment").List()
		question.Type("reactions", "Reaction").List()
	})

	schema.Type("Review", func(review *dqlx.TypeBuilder) {
		review.Predicate(name)
		review.Predicate(description)
		review.Predicate(moderation)
		review.Int("rating")
		review.Type("product", "Product").Reverse()
		review.Type("author", "User").Reverse()
		review.Type("business", "User")
		review.Type("assets", "Asset").List()
	})

	schema.Type("Reaction", func(reaction *dqlx.TypeBuilder) {
		reaction.Predicate(name)
		//reaction.Type("author", "User").Reverse()
	})

	schema.Type("HashTag", func(hashtag *dqlx.TypeBuilder) {
		hashtag.Predicate(name)
	})

	fmt.Println("End Of Schema")

}

func init() {
	fmt.Println("init function here")
	var err error
	db, err = dqlx.Connect("localhost:9080")
	/*conn, err := grpc.Dial("localhost:9080", grpc.WithInsecure())

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
	`*/

	if err != nil {
		log.Fatal(err)
	}
}

func GetDB() dqlx.DB {
	return db
}
