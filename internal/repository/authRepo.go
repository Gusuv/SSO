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

				return fmt.Errorf("create user %w", ErrUserExist)
			}

			return err
		}

		err := a.txSetRole(ctx, tx, users.Id)
		if err != nil {
			return fmt.Errorf("set role error %w", ErrSetRoleError)
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

func (a *AuthRep) GetUserByEmail(ctx context.Context, email string) (*models.Users, *models.Roles, error) {

	user := models.Users{
		Email: email,
	}
	role := models.Roles{}
	userRoles := models.UsersRoles{}

	if err := a.db.WithContext(ctx).Where("email = ?", user.Email).First(&user).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, fmt.Errorf("user not found %w", ErrUserNotFound)
		}
		return nil, nil, fmt.Errorf("get user by email %w", err)
	}

	if err := a.db.WithContext(ctx).Where("user_id = ?", user.Id).First(&userRoles).Error; err != nil {
		return nil, nil, fmt.Errorf("get user roles %w", err)
	}
	if err := a.db.WithContext(ctx).Where("id = ?", userRoles.RoleId).First(&role).Error; err != nil {
		return nil, nil, fmt.Errorf("get role %w", err)
	}

	return &user, &role, nil
}
