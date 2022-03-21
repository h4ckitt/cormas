package models

import "time"

type Balance struct {
	Amount   float64  `json:"amount"`
	Currency Currency `json:"currency"`
	//This Relationship Should Be A Reversible One To User
	Status     int       `json:"status"`
	Moderation bool      `json:"moderation"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
