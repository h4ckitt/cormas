package models

type Business struct {
	Name string `json:"name"`
}

type Debate struct {
	Name string `json:"name"`
}

type Group struct {
	Name string `json:"name"`
}

type Invoice struct {
	Number string `json:"number"`
}

type Product struct {
	Name string `json:"name"`
}

type Post struct {
	Name string `json:"name"`
}

type Comment struct {
	Body string `json:"body"`
}

type Question struct {
	Name string `json:"name"`
}

type Answer struct {
	Body string `json:"body"`
}

type Contact struct {
	Name string `json:"name"`
}
