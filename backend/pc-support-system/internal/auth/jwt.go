package auth

import (
	"errors"
	"time"

	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
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

func ValidateToken(tokenString, secret string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			return []byte(secret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func GenerateRefreshToken(
	user *models.User,
	secret string,
) (string, time.Time, error) {

	expiresAt := time.Now().Add(30 * 24 * time.Hour)

	claims := jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"type":    "refresh",
		"exp":     expiresAt.Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

func ValidateRefreshToken(
	tokenString string,
	secret string,
) (bson.ObjectID, error) {

	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}

			return []byte(secret), nil
		},
	)

	if err != nil {
		return bson.ObjectID{}, err
	}

	if !token.Valid {
		return bson.ObjectID{}, errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return bson.ObjectID{}, errors.New("invalid token claims")
	}

	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return bson.ObjectID{}, errors.New("invalid token type")
	}

	userIDString, ok := claims["user_id"].(string)
	if !ok {
		return bson.ObjectID{}, errors.New("invalid user id")
	}

	userID, err := bson.ObjectIDFromHex(userIDString)
	if err != nil {
		return bson.ObjectID{}, errors.New("invalid user id")
	}

	return userID, nil
}
