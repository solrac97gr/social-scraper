package middleware

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/solrac97gr/telegram-followers-checker/database"
)

type JWTConfig struct {
	Secret string
}

type AuthMiddleware struct {
	config *JWTConfig
}

type Claims struct {
	UserID       string                `json:"user_id"`
	Email        string                `json:"email"`
	Role         database.Role         `json:"role"`
	Subscription database.Subscription `json:"subscription"`
	jwt.RegisteredClaims
}

func NewAuthMiddleware(config *JWTConfig) (*AuthMiddleware, error) {
	if config == nil || config.Secret == "" {
		return nil, fmt.Errorf("JWT config and secret are required")
	}
	return &AuthMiddleware{
		config: config,
	}, nil
}

func (m *AuthMiddleware) WithJWT() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header required",
			})
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format",
			})
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token required",
			})
		}

		// Parse and validate the token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(m.config.Secret), nil
		})

		if err != nil {
			log.Printf("JWT parsing error: %v", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		if !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		// Store user information in the context
		c.Locals("userID", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)
		c.Locals("subscription", claims.Subscription)

		log.Printf("JWT authenticated user: %s (%s)", claims.Email, claims.UserID)
		return c.Next()
	}
}
