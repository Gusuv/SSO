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

		err := a.setRole(ctx, tx, users.Id)
		if err != nil {
			return fmt.Errorf("set role error %w", ErrSetRoleError)
		}
		userID = users.Id
		return nil

	})
	return userID, err
}

func (a *AuthRep) setRole(ctx context.Context, tx *gorm.DB, userId int64) error {

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

func (a *AuthRep) GetUserWithRole(ctx context.Context, email string) (*models.Users, []string, error) {
	type result struct {
		models.Users
		Role string
	}
	var res result

	if err := a.db.WithContext(ctx).Table("users").
		Select("users.id, users.username, users.password_hash, roles.role").
		Joins("JOIN users_roles ur ON ur.user_id = users.id").
		Joins("JOIN roles ON roles.id = ur.role_id").
		Where("users.email = ?", email).First(&res).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, fmt.Errorf("user not found: %w", ErrUserNotFound)
		}
		return nil, nil, fmt.Errorf("get user with role: %w", err)
	}

	roleList := []string{res.Role}
	return &res.Users, roleList, nil
}

func (a *AuthRep) AddRefreshToken(ctx context.Context, jwt *models.JWT) error {
	refreshToken := models.RefreshTokens{
		UserId:    jwt.UserId,
		TokenHash: jwt.RefreshHash,
		ExpiresAt: &jwt.ExpiresAt,
	}
	if err := a.db.WithContext(ctx).Create(&refreshToken).Error; err != nil {

		return fmt.Errorf("add refresh token: %w", err)
	}
	return nil
}
