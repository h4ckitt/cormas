package models

type User struct {
	UID        uint64    `json:"uid"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Username   string    `json:"username"`
	Password   string    `json:"-"`
	Avatar     string    `json:"avatar,omitempty"`
	Cover      string    `json:"cover,omitempty"`
	Addresses  []Address `json:"address"`
	Moderation int       `json:"moderation"`
	IsBusiness int       `json:"is_business"`
	Verified   bool      `json:"verified"`
	Assets     Asset     `json:"assets"`
	Premium    int       `json:"premium"`
	Amount     float64   `json:"amount,omitempty"`
	LastIP     string    `json:"last_ip"`
	Currency   Currency  `json:"currency"`
	Bank       Balance   `json:"balance"`
}

type Address struct {
}

type Balance struct {
}

type Asset struct {
}

type Currency struct {
}
