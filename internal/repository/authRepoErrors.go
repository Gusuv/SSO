package repository

import "errors"

var (
	UserExist    = errors.New("user already exist")
	SetRoleError = errors.New("can`t set user role")
)
