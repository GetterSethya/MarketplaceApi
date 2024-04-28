package entities

import (
	"database/sql"
	"time"
)

type Transaction struct {
	ID        string  `json:"id"`
	Status    string  `json:"status"` // enum (menunggu, diterima seller, dalam pengiriman, diterima)
	ProductId string  `json:"productId"`
	BuyerId   string  `json:"buyerId"`
	SellerId  string  `json:"sellerId"`
	Total     float64 `json:"total"`
	Quantity  int     `json:"quantity"`
	Notes     string  `json:"notes"`

	CreatedAt time.Time    `json:"-"`
	UpdatedAt time.Time    `json:"-"`
	DeletedAt sql.NullTime `json:"-"`
}

type TransactionMinimal struct {
	ID       string  `json:"id"`
	Status   string  `json:"status"` // enum (menunggu, diterima seller, dalam pengiriman, diterima)
	Total    float64 `json:"total"`
	Quantity int     `json:"quantity"`
	Notes    string  `json:"notes"`

	CreatedAt time.Time    `json:"-"`
	UpdatedAt time.Time    `json:"-"`
	DeletedAt sql.NullTime `json:"-"`
}
