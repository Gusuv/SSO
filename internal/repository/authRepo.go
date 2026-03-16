package repository

import (
	"context"
	"errors"
	"fmt"
	"main/internal/models"

	"gorm.io/gorm"
)

type AuthRep struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *AuthRep {
	return &AuthRep{db: db}
}

const (
	roleAdmin string = "Admin"
	roleUser  string = "User"
)

func (a *AuthRep) TxCreateUser(ctx context.Context, username, email, passwordHash string) (int64, error) {
	var userID int64

	err := a.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		users := models.Users{
			Username:     username,
			Email:        email,
			PasswordHash: passwordHash,
		}

		if err := tx.Create(&users).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {

				return fmt.Errorf("create user %w", UserExist)
			}

			return err
		}

		err := a.txSetRole(ctx, tx, users.Id)
		if err != nil {
			return fmt.Errorf("set role error %w", SetRoleError)
		}
		userID = users.Id
		return nil

	})
	return userID, err
}

func (a *AuthRep) txSetRole(ctx context.Context, tx *gorm.DB, userId int64) error {

	roles := models.Roles{}

	if err := tx.WithContext(ctx).Where("role = ?", roleUser).First(&roles).Error; err != nil {
		return err
	}

	usersRoles := models.UsersRoles{
		UserId: userId,
		RoleId: roles.Id,
	}
	if err := tx.Create(&usersRoles).Error; err != nil {
		return err
	}
	return nil
}
