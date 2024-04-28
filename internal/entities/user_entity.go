package entities

import (
	"database/sql"
	"time"
)

type User struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Username     string `json:"username"`
	HashPassword string `json:"password"`

	CreatedAt time.Time    `json:"-"`
	UpdatedAt time.Time    `json:"-"`
	DeletedAt sql.NullTime `json:"-"`
}

type UserMinimal struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`

	CreatedAt time.Time    `json:"-"`
	UpdatedAt time.Time    `json:"-"`
	DeletedAt sql.NullTime `json:"-"`
}
