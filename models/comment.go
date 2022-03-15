package models

import "time"

type Comment struct {
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Moderation  bool        `json:"moderation"`
	Author      User        `json:"author"`
	Address     Address     `json:"address,omitempty"`
	Reply       *[]Comment  `json:"reply,omitempty"`
	Reaction    *[]Reaction `json:"reaction"`
	Post		Post        `json:"post"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAT   time.Time   `json:"updated_at"`
}
