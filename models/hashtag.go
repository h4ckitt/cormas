package models

type HashTag struct {
	UID  string `json:"uid"`
	Name string `json:"name,omitempty"`
	Type string `json:"dgraph.type,omitempty"`
}
