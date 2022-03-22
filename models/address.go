package models

type Address struct {
	Name       string `json:"name"`
	Address1   string `json:"address1"`
	Address2   string `json:"address2,omitempty"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	Status     int    `json:"status"`
	Moderation bool   `json:"moderation"`
	Type       string `json:"dgraph.type,omitempty"`
	//Relationship To User Table Should Be Reverse  So We Can Find Out Who Lives At This Address
	//Relationship To Post Should Be Reverse So We Can Find Out Where This Post Was Made
	//Relationship To Comment Should Be Reverse So We Can Find Out Where This Comment Was Made
	//Relationship To Order Should Be Reverse So We Can Find Out Which Order Was Made From This Address
	//Relationship To Review Should Be Reverse So We Can Find Out Which Review Was Posted At This Address
	//Relationship To QNA Should Be Reverse So We Can Find Out Which Questions Were Asked From This Address
}
