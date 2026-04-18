package models

import "time"

type RefreshTokens struct {
	Id        int64 `gorm:"primary key"`
	UserId    int64
	TokenHash string
	CreatedAt *time.Time `gorm:"autoCreateTime"`
	ExpiresAt *time.Time `gorm:"autoCreateTime"`
	Revoked   bool       `gorm:"default:false"`
}
