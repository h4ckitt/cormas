package models

type Asset struct {
	UID        string      `json:"uid,omitempty"`
	Name       string      `json:"name"`
	Image      string      `json:"image,omitempty"`
	Video      string      `json:"video,omitempty"`
	Document   string      `json:"document,omitempty"`
	Zip        string      `json:"zip,omitempty"`
	UpdatedAt  string      `json:"updated_at"`
	CreatedAt  string      `json:"created_at"`
	EntityUID  string      `json:"entity_uid,omitempty"`
	Author     interface{} `json:"author"`
	Moderation bool        `json:"moderation"`
	Type       string      `json:"dgraph.type,omitempty"`
	//Reverse From Post Or Product Or Comment Or Review Or Question
}
