package hash

import (
	"golang.org/x/crypto/bcrypt"
)

func HashCompare(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}

	return true
}
