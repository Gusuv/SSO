package service

import "errors"

var (
	ErrUserAlreadyExist = errors.New("user already exist")
	ErrUserCreating     = errors.New("error when creating a user")
	ErrPasswordHashing  = errors.New("error when hashing password")
	ErrUserNotFound     = errors.New("user not found")
	ErrInvalidPassword  = errors.New("invalid password")
)
