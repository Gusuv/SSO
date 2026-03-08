package repository

import (
	"context"
	"main/internal/models"

	"gorm.io/gorm"
)

type AuthRep struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *AuthRep {

	return &AuthRep{db: db}
}

func (a *AuthRep) CreateUser(ctx context.Context, username, email, passwordHash string) error {
	return a.db.Create(&models.Users{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
	}).Error
}
