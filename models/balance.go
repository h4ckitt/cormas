package models

type Balance struct {
	Name       string      `json:"name"`
	Amount     float64     `json:"amount"`
	Currency   Currency    `json:"currency"`
	User       interface{} `json:"user"`
	Status     int         `json:"status"`
	Moderation bool        `json:"moderation"`
	CreatedAt  string   `json:"created_at"`
	UpdatedAt  string   `json:"updated_at"`
	Type       string      `json:"dgraph.type,omitempty"`
}
