package models

type Bank struct {
	UID      string   `json:"uid,omitempty"`
	Name     string   `json:"name"`
	Amount   float64  `json:"amount"`
	Currency Currency `json:"currency,omitempty"`
	//User       interface{} `json:"user"`
	Status     int    `json:"status"`
	Moderation bool   `json:"moderation"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	Type       string `json:"dgraph.type,omitempty"`
}
