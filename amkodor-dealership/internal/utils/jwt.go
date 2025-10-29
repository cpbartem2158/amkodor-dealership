package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims структура claims для JWT токена
type JWTClaims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// RefreshJWT обновляет JWT токен
func RefreshJWT(tokenString string, secret string, expireHours int) (string, error) {
	claims, err := ValidateJWTClaims(tokenString, secret)
	if err != nil {
		return "", err
	}

	// Создаем новый токен с теми же данными
	return GenerateJWT(claims.UserID, claims.Email, secret, expireHours)
}

// ExtractUserIDFromToken извлекает UserID из токена
func ExtractUserIDFromToken(tokenString string, secret string) (int, error) {
	claims, err := ValidateJWTClaims(tokenString, secret)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}

// ValidateJWTClaims валидирует JWT токен и возвращает claims
func ValidateJWTClaims(tokenString string, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}