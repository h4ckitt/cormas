package models

import "time"

type User struct {
	UID         uint64     `json:"uid"`
	Name        string     `json:"name"`
	Email       string     `json:"email"`
	Username    string     `json:"username"`
	Password    string     `json:"password,omitempty"`
	Avatar      string     `json:"avatar,omitempty"`
	Cover       string     `json:"cover,omitempty"`
	Addresses   *[]Address `json:"address"`
	Moderation  bool       `json:"moderation"`
	IsBusiness  bool       `json:"is_business"`
	Verified    bool       `json:"verified"`
	Assets      *Asset     `json:"assets"` //This Should Be Reversed From Either Product Or Post
	Premium     int        `json:"premium"`
	Amount      float64    `json:"amount,omitempty"`
	LastIP      string     `json:"last_ip"`
	Currency    *Currency  `json:"currency"`
	Bank        *Balance   `json:"balance"`
	Order       *[]Order   `json:"order,omitempty"`
	Invoice     *[]Invoice `json:"invoice,omitempty"`
	Post        *[]Post    `json:"post,omitempty"`
	Review      *[]Review  `json:"review,omitempty"`
	Owner       *[]User    `json:"owner,omitempty"`
	Editor      *User      `json:"editor,omitempty"`
	Publisher   *User      `json:"publisher,omitempty"`
	Category    *Category  `json:"category,omitempty"`
	Sale        *[]Order   `json:"sale,omitempty"`
	SaleInvoice *[]Invoice `json:"sale_invoice,omitempty"`
	Privacy     int        `json:"privacy"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
