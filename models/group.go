package models

type Group struct {
	Name       string      `json:"name"`
	Icon       string      `json:"icon,omitempty"`
	About      string      `json:"about,omitempty"`
	Status     int         `json:"status"`
	Moderation int         `json:"moderation"`
	Owner      interface{} `json:"owner"`
	Editor     interface{} `json:"editor"`
	Publisher  interface{} `json:"publisher"`
	ParentUID  string      `json:"child,omitempty"`
	CreatedAt  string      `json:"created_at"`
	UpdatedAt  string      `json:"updated_at"`
	Type       string      `json:"dgraph.type,omitempty"`
}
