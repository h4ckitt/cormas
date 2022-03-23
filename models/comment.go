package models

type Comment struct {
	UID string `json:"uid,omitempty"`
	//Title       string      `json:"title"`
	Description string      `json:"description"`
	Moderation  bool        `json:"moderation"`
	Author      interface{} `json:"author"`
	Address     *Address    `json:"address,omitempty"`
	Reply       *[]Comment  `json:"reply,omitempty"`
	Reaction    *[]Reaction `json:"reaction,omitempty"`
	Post        *Post       `json:"post,omitempty"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
	EntityUID   string      `json:"entity_uid,omitempty"` //Should Be Either Of Pro
	Type        string      `json:"dgraph.type,omitempty"`
}
