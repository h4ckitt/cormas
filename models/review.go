package models

type Review struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Author      User     `json:"author"`
	Assets      []*Asset `json:"assets,omitempty"`
	Moderation  bool     `json:"moderation"`
	Business    *User    `json:"business,omitempty"`
	Product     *Product `json:"product,omitempty"`
	Rating      int      `json:"rating"`
}
