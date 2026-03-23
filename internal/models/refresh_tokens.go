package models

import "time"

type Refresh_tokens struct {
	Id        int64 `gorm:"primary key"`
	UserId    int64
	tokenHash string
	CreatedAt *time.Time `gorm:"autoCreateTime"`
	ExpiresAt *time.Time `gorm:"autoCreateTime"`
}
