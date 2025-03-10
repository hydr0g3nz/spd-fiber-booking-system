package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Logging middleware for request logging
func Logging() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Start timer
		start := time.Now()

		// Store the request path
		path := c.Path()

		// Process request
		err := c.Next()

		// Calculate response time
		responseTime := time.Since(start)

		// Log request details
		log.Printf(
			"[%s] %s - Status: %d - Response time: %v",
			c.Method(),
			path,
			c.Response().StatusCode(),
			responseTime,
		)

		return err
	}
}
