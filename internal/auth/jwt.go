package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// SecretKey используется для подписи токенов
var SecretKey = []byte(os.Getenv("JWT_SECRET"))

// GenerateToken создаёт JWT токен для пользователя
func GenerateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Токен действует 24 часа
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SecretKey)
}
