package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

type Hash struct {
	hmacSecret string
}

func NewHash(hmacSecret string) *Hash {
	return &Hash{hmacSecret: hmacSecret}
}

func (h *Hash) MakeHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 8)

	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (h *Hash) HashToken(refreshToken string) string {

	mac := hmac.New(sha256.New, []byte(h.hmacSecret))
	mac.Write([]byte(refreshToken))

	tokenHash := hex.EncodeToString(mac.Sum(nil))

	return tokenHash
}

func (h *Hash) HashCompare(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}

	return true
}
