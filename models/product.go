package models

type Product struct {
	Name         string    `json:"name"`
	RegularPrice float64   `json:"regular_price"`
	SellingPrice float64   `json:"selling_price"`
	Currency     *Currency `json:"currency"`
	Category     *Category `json:"category"`
	//This Should Be A Reverse Relationship To Order Table
	Reviews               *[]Review `json:"review"`
	ProductType           int       `json:"type"`
	Supported             int       `json:"supported"`
	Downloadable          *[]Asset  `json:"downloadable,omitempty"`
	Thumbnail             *Asset    `json:"thumbnail,omitempty"`
	Gallery               *[]Asset  `json:"gallery"`
	Excerpt               string    `json:"excerpt"`
	Description           string    `json:"description"`
	TechnicalInformation  string    `json:"technical_information,omitempty"`
	AdditionalInformation string    `json:"additional_information,omitempty"`
	ProductInformation    string    `json:"product_information,omitempty"`
	ProductGuides         string    `json:"product_guides,omitempty"`
	Status                int       `json:"status"`
	Moderation            bool      `json:"moderation"`
	Address               *Address  `json:"address"`
	Owner                 *User     `json:"owner"`
	Type                  string    `json:"dgraph.type,omitempty"`
}
