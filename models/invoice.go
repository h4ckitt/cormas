package models

import "time"

type Invoice struct {
	Order     *[]Order  `json:"order,omitempty"`
	Buyer     *User     `json:"buyer"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Status    int       `json:"status"`
}
