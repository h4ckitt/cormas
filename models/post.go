package models

type Post struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Author      interface{} `json:"author"`
	Business    interface{} `json:"business"`
	Orders      *[]Order    `json:"orders"`
	Privacy     int         `json:"privacy"`
	Address     Address     `json:"address"`
	Moderation  bool        `json:"moderation"`
	Currency    Currency    `json:"currency"`
	Comments    *[]Comment  `json:"comments,omitempty"`
	Amount      float64     `json:"amount"`
	Reactions   *[]Reaction `json:"reactions,omitempty"`
	Assets      *[]Asset    `json:"assets,omitempty"`
	Tags        *[]HashTag  `json:"tags,omitempty"`
	Type        string      `json:"dgraph.type"`
}
