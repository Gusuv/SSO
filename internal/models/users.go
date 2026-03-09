package models

import (
	"time"
)

type Users struct {
	Id           int64      `gorm:"primary key"`
	Username     string     `gorm:"not null"`
	Email        string     `gorm:"not null"`
	PasswordHash string     `gorm:"not null"`
	CreatedAt    *time.Time `gorm:"autoCreateTime"`
	UpdatedAt    *time.Time `gorm:"autoCreateTime"`
}
