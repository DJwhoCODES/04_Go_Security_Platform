package security

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secret string
	issuer string
}

type Claims struct {
	UserID string   `json:"user_id"`
	Roles  []string `json:"roles"`

	jwt.RegisteredClaims
}

func NewJWTManager(secret string, issuer string) *JWTManager {

	return &JWTManager{
		secret: secret,
		issuer: issuer,
	}
}

func (j *JWTManager) GenerateAccessToken(
	userID string,
	roles []string,
	expiry time.Duration,
) (string, error) {

	claims := Claims{
		UserID: userID,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(j.secret))
}

func (j *JWTManager) Verify(tokenStr string) (*Claims, error) {
	claims := &Claims{}

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		secretKey := []byte(j.secret)
		return secretKey, nil
	}

	parsedToken, err := jwt.ParseWithClaims(
		tokenStr,
		claims,
		keyFunc,
	)

	if err != nil {
		return nil, err
	}

	extractedClaims, ok := parsedToken.Claims.(*Claims)

	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return extractedClaims, nil
}

func GenerateRefreshToken() (string, error) {

	return GenerateRandomToken(32)
}
