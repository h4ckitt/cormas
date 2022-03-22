package models

type Currency struct {
	Name      string `json:"name"`
	Icon      string `json:"icon,omitempty"`
	Value     string `json:"value"`
	Status    int    `json:"status"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	Type      string `json:"dgraph.type"`
}
