package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// Auth middleware for authentication
func Auth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// This is a mock authentication middleware
		// In a real application, you would validate tokens, check permissions, etc.

		// For demo purposes, we'll just check for a mock API key
		apiKey := c.Get("X-API-Key")
		// if apiKey == "" {
		// 	// If no API key provided, let it pass (for demo/testing)
		// 	return c.Next()
		// }

		// For demo, we'll accept any API key that's at least 10 characters
		if len(apiKey) < 10 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid API Key",
			})
		}

		// Authentication successful
		return c.Next()
	}
}
