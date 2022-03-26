package models

type Review struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Author      interface{}  `json:"author"`
	Assets      []*Asset     `json:"assets,omitempty"`
	Moderation  bool         `json:"moderation"`
	Business    *interface{} `json:"business,omitempty"`
	Product     *interface{} `json:"product,omitempty"`
	Rating      int          `json:"rating"`
	CreatedAt   string       `json:"created_at"`
	UpdatedAt   string       `json:"updated_at"`
	Type        string       `json:"dgraph.type,omitempty"`
}
