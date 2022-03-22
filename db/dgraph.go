package db

import (
	"context"
	"fmt"
	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"google.golang.org/grpc"

	//"google.golang.org/grpc"
	"log"
)

var db *dgo.Dgraph

//var db dqlx.DB

/*func CreateSchema(drop bool) {
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

}*/

func DropAll() {
	if err := db.Alter(context.Background(), &api.Operation{DropAll: true}); err != nil {
		log.Println("Couldn't Wipe Database, Work With What You Have")
		return
	}

	log.Println("Wiped Database Successfully")
}

func CreateSchema() {
	op := api.Operation{}

	op.Schema = `
			name: string @index(exact) .
			moderation: int @index(int) .
			user_agent: string .
			privacy: int .
			amount: float .
			created_at: datetime .
			updated_at: datetime .
			status: int .
			suspended: bool .
			icon: string .
			description: string .
			username: string @index(hash) .
			email: string @index(hash) .
			is_business: bool . 
			verified: bool .
			password: password .
			avatar: string .
			cover: string .
			last_ip: string .
			balance: float .
			addresses: [uid] .
			currency: uid .
			bank: uid .
			orders: [uid] .
			invoices: [uid] .
			posts: [uid] .
			reviews: [uid] .
			owners: [uid] .
			editor: uid .
			publisher: uid .
			category: uid .
			sales: [uid] .
			sales_invoices: [uid] .
			author: uid @reverse . 
			business: uid .
			address: uid .
			comments: [uid] .
			reactions: [uid] .
			tags: [uid] .
			assets: [uid] .
			transaction_id: string .
			sender: uid .
			receiver: uid .
			products: [uid] .
			address1: string .
			address2: string .
			city: string .
			country: string .
			image: string .
			video: string .
			document: string .
			zip: string .
			post: uid @reverse .
			comment: uid @reverse .
			review: uid @reverse .
			value: string .
			child: [uid] @reverse .
			replies: [uid] .
			order: uid .
			buyer: uid .
			supported: int .
			type: int .
			excerpt: string .
			technical_information: string .
			additional_information: string .
			product_information: string .
			product_guides: string .
			regular_price: float .
			selling_price: float .
			owner: uid @reverse .
			thumbnail: string .
			downloadable: [uid] .
			gallery: [uid] .
			rating: int .
			product: uid @reverse .

			type User {
				name
				moderation
				privacy
				created_at
				updated_at
				amount
				username
				email
				is_business
				verified
				password
				avatar
				cover
				last_ip
				user_agent
				balance
				addresses: [Address]
				currency: Currency
				bank: Balance
				orders: [Order]
				invoices: [Invoice]
				posts: [Post]
				reviews: [Review]
				owners: [User]
				editor: User
				publisher: User
				category: Category
				sales: [Order]
				sales_invoices: [Invoice]
			}

			type Post {
				name
				description
				privacy
				moderation
				amount
				author: User
				business: User
				orders: [Order]
				address: Address
				currency: Currency
				comments: [Comment]
				reactions: [Reaction]
				tags: [HashTag]
				assets: [Asset]
			}

			type Order {
				status
				amount
				moderation
				user_agent
				created_at
				updated_at
				transaction_id
				business: User
				sender: User
				receiver: User
				currency: Currency
				posts: [Post]
				products: [Product]
			}

			type Address {
				name
				status
				moderation
				address1
				address2
				city
				country
			
			}

			type Asset {
				name
				created_at
				updated_at
				moderation
				image
				video
				document
				zip
				post: Post
				comment: Comment
				review: Review
			}

			type Currency {
				name
				icon
				status
				created_at
				updated_at
				value
			}

			type Balance {
				amount
				status
				moderation
				created_at
				updated_at
				currency: Currency
			}

			type Category {
				name
				icon
				moderation
				status
				child: [Category]
			}

			type Comment {
				description
				moderation
				created_at
				updated_at
				author: User
				post: Post
				reactions: [Reaction]
				replies: [Comment]
				address: Address
			}

			type Invoice {
				status
				created_at
				updated_at
				order: Order
				buyer: User
			}

			type Product {
				name
				description
				moderation
				status
				supported
				type
				excerpt
				technical_information
				additional_information
				product_information
				product_guides
				regular_price
				selling_price
				address: Address
				owner: User
				currency: Currency
				category: Category
				reviews: [Review]
				thumbnail
				downloadable: [Asset]
				gallery: [Asset]
			}

			type Question {
				name
				description
				status
				moderation
				privacy
				created_at
				updated_at
				tags: [HashTag]
				comments: [Comment]
				reactions: [Reaction]
			}

			type Review {
				name
				description
				moderation
				rating
				product: Product
				author: User
				business: User
				assets: [Asset]
			}

			type Reaction {
				name
			}

			type HashTag {
				name
			}
	`

	if err := db.Alter(context.Background(), &op); err != nil {
		log.Println("Couldn't Create The Schema, Work With What You Have: ", err)
	}

}

func init() {
	fmt.Println("init function here")
	var err error
	//db, err = dqlx.Connect("localhost:9080")
	conn, err := grpc.Dial("localhost:9080", grpc.WithInsecure())

	if err != nil {
		log.Fatal(err)
	}

	db = dgo.NewDgraphClient(api.NewDgraphClient(conn))

}

func GetDB() *dgo.Dgraph {
	return db
}
