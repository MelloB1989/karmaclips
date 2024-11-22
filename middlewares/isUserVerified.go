package middlewares

import (
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"

	"karmaclips/config"
)

// Define a custom struct to hold the JWT claims
type Claims struct {
	UserID string `json:"uid"`
	Phone  string `json:"phone"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

// JWTMiddleware is the middleware function to verify JWT tokens
func IsUserVerified(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Missing or malformed JWT")
	}

	// Extract the token from the Bearer string
	tokenStr := strings.TrimSpace(strings.Replace(authHeader, "Bearer", "", 1))
	if tokenStr == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Missing or malformed JWT")
	}

	// Parse the JWT token
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			// Make sure the token's algorithm is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Unexpected signing method")
			}
			return []byte(config.NewConfig().JWTSecret), nil
		},
	)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid or expired JWT")
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// Store the claims in the context's locals
		// fmt.Println(claims.UserID)
		c.Locals("uid", claims.UserID)
		c.Locals("phone", claims.Phone)
		c.Locals("email", claims.Email)
		c.Locals("exp", time.Unix(claims.ExpiresAt, 0))
		// Continue with the next handler
		return c.Next()
	}

	return c.Status(fiber.StatusUnauthorized).SendString("Invalid or expired JWT")
}
