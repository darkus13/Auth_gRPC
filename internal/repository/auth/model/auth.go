package model

import (
	"database/sql"
	"time"
)

type Auth struct {
	ID        int64        `db:"id"`
	Info      AuthInfo     `db:""`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

type AuthInfo struct {
	Name        string `db:"name"`
	Email       string `db:"email"`
	Password    string `db:"password"`
	PassConfirm string `db:"password_confirm"`
	Role        int32  `db:"role"`
}
