package service

import "errors"

var (
	ErrUserAlreadyExist    = errors.New("user already exist")
	ErrUserCreating        = errors.New("error when creating a user")
	ErrPasswordHashing     = errors.New("error when hashing password")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrRefreshTokenExpired = errors.New("refresh token expired")
)
