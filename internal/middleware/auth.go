package middleware

import (
	"log/slog"
	"net/http"
	"strings"

	"user-rewards-api/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	UserIDKey = "user_id"
)

// AuthMiddleware middleware для проверки JWT токена
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			slog.Warn("Попытка доступа без токена авторизации", "path", c.Request.URL.Path)
			sendError(c, "токен авторизации отсутствует")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			slog.Warn("Некорректный формат токена", "path", c.Request.URL.Path)
			sendError(c, "некорректный формат токена")
			c.Abort()
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, domain.ErrInvalidUsername
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			slog.Warn("Ошибка парсинга JWT токена", "error", err, "path", c.Request.URL.Path)
			sendError(c, "невалидный токен: "+err.Error())
			c.Abort()
			return
		}

		if !token.Valid {
			slog.Warn("Невалидный JWT токен", "path", c.Request.URL.Path)
			sendError(c, "токен невалиден")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			slog.Warn("Некорректные claims в JWT токене", "path", c.Request.URL.Path)
			sendError(c, "некорректные claims токена")
			c.Abort()
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok || userID == "" {
			slog.Warn("user_id не найден в JWT токене", "path", c.Request.URL.Path)
			sendError(c, "user_id не найден в токене")
			c.Abort()
			return
		}

		c.Set(UserIDKey, userID)
		c.Next()
	}
}

// sendError отправляет ошибку в формате JSON
func sendError(c *gin.Context, message string) {
	response := map[string]string{
		"error": message,
	}
	c.JSON(http.StatusUnauthorized, response)
}
