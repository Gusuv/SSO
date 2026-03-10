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

func (a *AuthRep) CreateUser(ctx context.Context, username, email, passwordHash string) (error, int64) {
	user := models.Users{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
	}

	if err := a.db.WithContext(ctx).Create(&user).Error; err != nil {
		return err, 0
	}
	return nil, user.Id

}

func (a *AuthRep) SetRole(ctx context.Context, id int64) error {
	role := models.Roles{}
	a.db.Where("role = ?", "User").First(&role)
	return a.db.Create(&models.UsersRoles{
		UserId: id,
		RoleId: role.Id,
	}).Error
}
