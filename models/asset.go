package models

import "time"

type Asset struct {
	Name      string    `json:"name"`
	Image     string    `json:"image,omitempty"`
	Video     string    `json:"video,omitempty"`
	Document  string    `json:"document,omitempty"`
	Zip       string    `json:"zip,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
	//Reverse From Post Or Product Or Comment Or Review Or Question
	Moderation bool `json:"moderation"`
}
