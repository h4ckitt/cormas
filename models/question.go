package models

import "time"

type Question struct {
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Tags        *[]HashTag  `json:"tags,omitempty"`
	Status      int         `json:"status"`
	Moderation  bool        `json:"moderation"`
	Privacy     int         `json:"privacy"`
	Author      User        `json:"author"`
	Address     Address     `json:"address"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	Comments    *[]Comment  `json:"comments,omitempty"`
	Reactions   *[]Reaction `json:"reactions,omitempty"`
}
