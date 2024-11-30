package middlewares

import (
	"karmaclips/config"

	"github.com/gofiber/fiber/v2"
)

func IsKarmaAdmin(c *fiber.Ctx) error {
	authHeader := c.GetReqHeaders()
	// Check if the header exists and has at least one value
	if values, exists := authHeader["X-Karma-Admin-Auth"]; exists && len(values) > 0 {
		// Now compare the first value of the slice to the expected value
		if values[0] == config.NewConfig().AdminKey {
			c.Locals("admin", true)
			return c.Next()
		}
	}
	// If the header does not exist, is empty, or does not match "karma-admin", handle accordingly
	// For example, you might want to return an error or unauthorized response
	return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
}
