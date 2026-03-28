package security

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AccessClaims struct {
	UserId string
	RoleId string
	AppId  string
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
func (j *JWTService) GenerateToken(userId, appId int64, role string) (string, string, error) {

	claimss := AccessClaims{
		UserId: strconv.FormatInt(userId, 10),
		RoleId: role,
		AppId:  strconv.FormatInt(appId, 10),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.TokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		}}

	accesToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claimss).SignedString(j.Secret)
	if err != nil {
		return "", "", err
	}

	refreshToken := uuid.NewString()

	return accesToken, refreshToken, nil
}
