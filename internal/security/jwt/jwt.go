package security

import (
	"time"
)

type AccessClaims struct {
	UserId    string
	AppId     string
	ExpiresAt int64
}

type JWTService struct {
	Secret   string
	TokenTTL time.Duration
}

func NewToken(secret string, tokenTTL time.Duration) *JWTService {
	return &JWTService{
		Secret:   secret,
		TokenTTL: tokenTTL,
	}
}
func (t *JWTService) GenerateToken(userId string) (string, string, error) {

	return "", "", nil

}
