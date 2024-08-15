package models

import (
	"time"
)

type User struct {
	ID         string    `json:"id"`
	Account    string    `json:"account"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	UserInfo   UserInfo  `json:"user_info"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}

type UserInfo struct {
	Name    string `json:"name"`
	Country string `json:"country"`
	Points  int64  `json:"points"`
	Hp      int32  `json:"hp"`
}
