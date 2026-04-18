package repository

import "errors"

var (
	ErrUserExist            = errors.New("user already exist")
	ErrSetRoleError         = errors.New("can`t set user role")
	ErrUserNotFound         = errors.New("email not found")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
)
