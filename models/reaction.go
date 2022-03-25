package models

type Reaction struct {
	UID       string      `json:"uid,omitempty"`
	Name      string      `json:"name"`
	EntityUID string      `json:"entity_uid"`
	Owner     interface{} `json:"owner"`
	CreatedAt string      `json:"created_at"`
	UpdatedAt string      `json:"updated_at"`
	Type      string      `json:"dgraph.type,omitempty"`
}
