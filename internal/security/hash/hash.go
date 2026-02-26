package hash

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func PasswordHash(password string) (error, string) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return errors.New("Password is not hashed"), ""
	}
	return nil, string(hash)
}
