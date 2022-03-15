package models

import "time"

type User struct {
	UID          uint64     `json:"uid"`
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	Username     string     `json:"username"`
	Password     string     `json:"password,omitempty"`
	Avatar       string     `json:"avatar,omitempty"`
	Cover        string     `json:"cover,omitempty"`
	Addresses    *[]Address `json:"address"`
	Moderation   bool       `json:"moderation"`
	IsBusiness   bool       `json:"is_business"`
	Verified     bool       `json:"verified"`
	Assets       *Asset     `json:"assets"` //This Should Be Reversed From Either Product Or Post
	Premium      int        `json:"premium"`
	Amount       float64    `json:"amount,omitempty"`
	LastIP       string     `json:"last_ip"`
	Currency     *Currency  `json:"currency"`
	Bank         *Balance   `json:"balance"`
	Orders       *[]Order   `json:"orders,omitempty"`
	Invoices     *[]Invoice `json:"invoices,omitempty"`
	Posts        *[]Post    `json:"posts,omitempty"`
	Reviews      *[]Review  `json:"reviews,omitempty"`
	Owners       *[]User    `json:"owners,omitempty"`
	Editor       *User      `json:"editor,omitempty"`
	Publisher    *User      `json:"publisher,omitempty"`
	Category     *Category  `json:"category,omitempty"`
	Sales        *[]Order   `json:"sales,omitempty"`
	SaleInvoices *[]Invoice `json:"sale_invoices,omitempty"`
	Privacy      int        `json:"privacy"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
