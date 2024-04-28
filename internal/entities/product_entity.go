package entities

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type Product struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Price          float64        `json:"price"`
	ImageUrl       string         `json:"imageUrl"`
	Stock          int            `json:"stock"`
	Condition      string         `json:"condition"`
	Tags           pq.StringArray `json:"tags" db:"tags"`
	IsPurchaseable bool           `json:"isPurchaseable"`
	SellerId       string         `json:"sellerId"`
	Descriptions   string         `json:"descriptions"`

	CreatedAt time.Time    `json:"-"`
	UpdatedAt time.Time    `json:"-"`
	DeletedAt sql.NullTime `json:"-"`
}

type ProductMinimal struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Price          float64        `json:"price"`
	ImageUrl       string         `json:"imageUrl"`
	Stock          int            `json:"stock"`
	Condition      string         `json:"condition"`
	Tags           pq.StringArray `json:"tags" db:"tags"`
	IsPurchaseable bool           `json:"isPurchaseable"`
	Descriptions   string         `json:"descriptions"`

	CreatedAt time.Time    `json:"-"`
	UpdatedAt time.Time    `json:"-"`
	DeletedAt sql.NullTime `json:"-"`
}
