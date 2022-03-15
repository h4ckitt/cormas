package models

type Post struct {
	Name       string     `json:"name"`
	Message    string     `json:"message"`
	Author     *User      `json:"author"`
	Business   *User      `json:"business"`
	Order      *Order     `json:"order"`
	Privacy    int        `json:"privacy"`
	Address    *[]Address `json:"address"`
	Moderation bool       `json:"moderation"`
	Currency   *Currency  `json:"currency"`
	Amount     float64    `json:"amount"`
	Reaction   *Reaction  `json:"reaction"`
	Assets     *[]Asset   `json:"asset"`
	Tags       *[]HashTag `json:"tags"`
}

type Reaction struct {
}

type HashTag struct {
}
