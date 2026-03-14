package validation

import "errors"

var (
	ShortPassword    = errors.New("password is too short")
	LongPassword     = errors.New("password is too long")
	RequiredEmail    = errors.New("email is required")
	RequiredPassword = errors.New("password is required")
	RequiredUsername = errors.New("username is required")
	RequiredAppId    = errors.New("app id is required")
	InvalidEmail     = errors.New("email is invalid")
	InvalidUsername  = errors.New("username is invalid")
	ShortUsername    = errors.New("username is too short")
	LongUsername     = errors.New("username is too long")
)
