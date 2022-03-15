package models

type Category struct {
	Name       string `json:"name"`
	Icon       string `json:"icon,omitempty"`
	Moderation bool   `json:"moderation"`
	Status     int    `json:"status"`
	//Relationship To Product Should Be Reverse
	//Relationship To Author/Business Should Be Reverse
	Child *[]Category `json:"child"`
}
