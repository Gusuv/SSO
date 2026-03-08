package hash

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func MakeHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", errors.New("Password is not hashed")
	}
	return string(hash), nil
}
