package models

import "time"

type Order struct {
	Post    *[]Post    `json:"post,omitempty"`
	Product *[]Product `json:"product,omitempty"`
	//The Relationship To Invoice should be reverse
	TransactionID string    `json:"transaction_id"`
	Sender        *User     `json:"sender,omitempty"`
	Business      *User     `json:"business,omitempty"`
	Receiver      *User     `json:"receiver,omitempty"`
	Status        int       `json:"status"`
	Moderation    bool      `json:"moderation"`
	Amount        float64   `json:"amount"`
	Currency      Currency  `json:"currency"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	UserAgent     string    `json:"user_agent,omitempty"`
}
