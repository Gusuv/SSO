package security

import (
	"main/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type accessClaims struct {
	UserId int64
	Role   string
	AppId  int64
	jwt.RegisteredClaims
}

type JWTService struct {
	Secret   []byte
	TokenTTL time.Duration
}

func NewToken(secret string, tokenTTL time.Duration) *JWTService {
	return &JWTService{
		Secret:   []byte(secret),
		TokenTTL: tokenTTL,
	}
}

func (j *JWTService) GenerateToken(userId, appId int64, role string) (*models.JWT, error) {

	claims := accessClaims{
		UserId: userId,
		Role:   role,
		AppId:  appId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.TokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		}}

	accesToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(j.Secret)
	if err != nil {
		return nil, err
	}

	refreshToken := uuid.NewString()

	jwt := models.JWT{
		AccessToken:  accesToken,
		RefreshToken: refreshToken,
		UserId:       userId,
		ExpiresAt:    time.Now().Add(j.TokenTTL),
	}

	return &jwt, nil
}
