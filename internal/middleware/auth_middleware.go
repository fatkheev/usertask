package middleware

import (
	"net/http"
	"strings"
	"usertask/internal/auth"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware проверяет JWT-токен перед выполнением запроса.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем заголовок Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "токен отсутствует"})
			c.Abort()
			return
		}

		// Ожидаемый формат: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "неверный формат токена"})
			c.Abort()
			return
		}

		tokenStr := parts[1]

		// Разбираем и проверяем токен
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return auth.SecretKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "недействительный токен"})
			c.Abort()
			return
		}

		// Извлекаем user_id из токена
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "не удалось получить данные пользователя"})
			c.Abort()
			return
		}

		// Добавляем user_id в контекст Gin
		c.Set("user_id", int(claims["user_id"].(float64)))

		// Передаём управление дальше
		c.Next()
	}
}
