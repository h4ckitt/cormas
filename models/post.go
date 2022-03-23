package models

type Post struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Author      interface{} `json:"author"`
	Business    interface{} `json:"business,omitempty"`
	Orders      *[]Order    `json:"orders,omitempty"`
	Privacy     int         `json:"privacy"`
	Address     Address     `json:"address,omitempty"`
	Moderation  bool        `json:"moderation"`
	Currency    *Currency   `json:"currency,omitempty"`
	Comments    *[]Comment  `json:"comments,omitempty"`
	Amount      float64     `json:"amount,omitempty"`
	Reactions   *[]Reaction `json:"reactions,omitempty"`
	Assets      *[]Asset    `json:"assets,omitempty"`
	Tags        *[]HashTag  `json:"tags,omitempty"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
	Type        string      `json:"dgraph.type,omitempty"`
}
