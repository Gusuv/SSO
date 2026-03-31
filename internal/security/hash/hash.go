package hash

import (
	"golang.org/x/crypto/bcrypt"
)

func MakeHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 8)

	if err != nil {
		return "", err
	}
	return string(hash), nil
}
