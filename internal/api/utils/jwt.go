package utils

import (
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

var jwtSecret = []byte(viper.GetString("jwt.secret"))

func GenerateJWT(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
	})

	return token.SignedString(jwtSecret)
}

func ExtractUserIDFromJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, exists := claims["user_id"].(string); exists {
			return userID, nil
		}
	}

	return "", errors.New("invalid token")
}

func GetUserIDFromHeader(authHeader string) (string, error) {
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("missing or invalid auth header")
	}
	return ExtractUserIDFromJWT(strings.TrimPrefix(authHeader, "Bearer "))
}
