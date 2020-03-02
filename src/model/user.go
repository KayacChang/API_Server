package model

import (
	"time"
)

type User struct {
	ID       string `json:"id" db:"user_id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
