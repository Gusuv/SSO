package security

import (
	"main/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type accessClaims struct {
	UserId int64
	Role   []string
	jwt.RegisteredClaims
}

type JWTService struct {
	Secret   []byte
	TokenTTL config.TokensTTL
}

func NewToken(secret string, tokenTTL config.TokensTTL) *JWTService {
	return &JWTService{
		Secret:   []byte(secret),
		TokenTTL: tokenTTL,
	}
}

func (j *JWTService) GenerateTokens(userId int64, role []string) (*Tokens, error) {
	now := time.Now()

	accessToken, err := j.GenerateAccessToken(userId, role)
	if err != nil {
		return nil, err
	}

	refreshToken := uuid.NewString()

	jwt := Tokens{
		AccessToken:      accessToken.AccessToken,
		AccessExpiresAt:  accessToken.AccessExpiresAt,
		RefreshExpiresAt: now.Add(j.TokenTTL.Refresh),
		UserId:           userId,
		RefreshToken:     refreshToken,
	}
	return &jwt, nil
}

func (j *JWTService) GenerateAccessToken(userID int64, role []string) (*AccessToken, error) {
	now := time.Now()
	claims := accessClaims{
		UserId: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(j.TokenTTL.Access)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(j.Secret)
	if err != nil {
		return nil, err
	}

	return &AccessToken{
		AccessToken:     accessToken,
		AccessExpiresAt: now.Add(j.TokenTTL.Access),
	}, nil
}

type AccessToken struct {
	AccessToken     string
	AccessExpiresAt time.Time
}

type Tokens struct {
	AccessToken      string
	AccessExpiresAt  time.Time
	RefreshExpiresAt time.Time
	UserId           int64
	RefreshToken     string
}
