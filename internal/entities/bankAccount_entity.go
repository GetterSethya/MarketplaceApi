package entities

import (
	"database/sql"
	"time"
)

type BankAccount struct {
	Id            string `json:"bankAccountId"`
	BankName      string `json:"bankName"`
	AccountName   string `json:"accountName"`
	AccountNumber int64  `json:"accountNumber"`
	SellerId      string `json:"sellerId"`

	CreatedAt time.Time    `json:"-"`
	UpdatedAt time.Time    `json:"-"`
	DeletedAt sql.NullTime `json:"-"`
}
