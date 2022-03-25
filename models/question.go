package models

import "time"

type Question struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Tags        *[]interface{} `json:"tags,omitempty"`
	Status      int            `json:"status"`
	Moderation  bool           `json:"moderation"`
	Privacy     int            `json:"privacy"`
	Author      interface{}    `json:"author"`
	Address     Address        `json:"address"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Comments    *[]Comment     `json:"comments,omitempty"`
	Reactions   *[]Reaction    `json:"reactions,omitempty"`
	Type        string         `json:"dgraph.type,omitempty"`
}
