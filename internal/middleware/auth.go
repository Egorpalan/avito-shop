package middleware

import (
	"github.com/Egorpalan/avito-shop/config"
	"github.com/Egorpalan/avito-shop/pkg/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем заголовок Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Проверяем, что заголовок начинается с "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		// Извлекаем токен (убираем "Bearer ")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Загружаем конфигурацию
		cfg := config.LoadConfig(".env.example")

		// Парсим токен
		claims, err := jwt.ParseJWT(tokenString, cfg.JWTSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Сохраняем имя пользователя в контексте
		c.Set("username", claims.Username)
		c.Next()
	}
}
