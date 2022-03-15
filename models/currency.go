package models

import "time"

type Currency struct {
	Name      string    `json:"name"`
	Icon      string    `json:"icon,omitempty"`
	Value     string    `json:"value"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
