package auth

import (
	"time"

	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(
	user *models.User,
	secret string,
	expiryHours int,
) (string, error) {

	claims := Claims{
		UserID: user.ID.Hex(),
		Role:   string(user.Role),

		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiryHours) * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}
