package repository

import (
	"context"
	"errors"
	"fmt"
	"main/internal/models"
	security "main/internal/security/jwt"

	"gorm.io/gorm"
)

type AuthRep struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *AuthRep {
	return &AuthRep{db: db}
}

const roleUser string = "User"

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
	var res []result

	if err := a.db.WithContext(ctx).Table("users").
		Select("users.id, users.username, users.password_hash, roles.role").
		Joins("JOIN users_roles ur ON ur.user_id = users.id").
		Joins("JOIN roles ON roles.id = ur.role_id").
		Where("users.email = ?", email).Find(&res).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, fmt.Errorf("user not found: %w", ErrUserNotFound)
		}
		if errors.Is(err, gorm.ErrInvalidValueOfLength) {
			return nil, nil, fmt.Errorf("invalid value of length: %w", err)
		}
		return nil, nil, fmt.Errorf("get user with role: %w", err)
	}

	roleList := make([]string, 0, len(res))
	for _, r := range res {
		roleList = append(roleList, r.Role)
	}
	if len(res) == 0 {
		return nil, nil, fmt.Errorf("user not found: %w", ErrUserNotFound)
	}

	user := &res[0].Users

	return user, roleList, nil
}

func (a *AuthRep) AddRefreshToken(ctx context.Context, jwt *security.Tokens, refreshHash string) error {
	refreshToken := models.RefreshTokens{
		UserId:    jwt.UserId,
		TokenHash: refreshHash,
		ExpiresAt: &jwt.RefreshExpiresAt,
	}
	if err := a.db.WithContext(ctx).Create(&refreshToken).Error; err != nil {

		return fmt.Errorf("add refresh token: %w", err)
	}
	return nil
}

func (a *AuthRep) CheckRefreshToken(ctx context.Context, tokenHash string) (*models.RefreshTokens, error) {
	var refToken models.RefreshTokens

	if err := a.db.Where("token_hash = ?", tokenHash).First(&refToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("refresh token not found: %w", ErrRefreshTokenNotFound)
		}
		return nil, fmt.Errorf("check refresh token: %w", err)
	}

	return &refToken, nil
}

func (a *AuthRep) GetUserRolesById(ctx context.Context, userId int64) ([]string, error) {
	type result struct {
		Role string
	}
	var role []result

	if err := a.db.WithContext(ctx).Table("users_roles AS ur").
		Select("r.role").
		Joins("JOIN roles r ON r.id = ur.role_id").
		Where("ur.user_id = ?", userId).
		Scan(&role).Error; err != nil {
		return nil, fmt.Errorf("get user roles: %w", err)
	}

	roleList := make([]string, 0, len(role))
	for _, r := range role {
		roleList = append(roleList, r.Role)
	}
	return roleList, nil
}

func (a *AuthRep) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	ref := models.RefreshTokens{
		Revoked: true,
	}
	if err := a.db.WithContext(ctx).Model(&ref).Where("token_hash = ?", tokenHash).Update("revoked", true).Error; err != nil {
		return fmt.Errorf("revoke refresh token: %w", err)
	}
	return nil
}
