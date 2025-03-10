package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/hydr0g3nz/spd-fiber-booking-system/handler"
	"github.com/hydr0g3nz/spd-fiber-booking-system/repository"
	"github.com/hydr0g3nz/spd-fiber-booking-system/router"
	"github.com/hydr0g3nz/spd-fiber-booking-system/usecase"
	"github.com/hydr0g3nz/spd-fiber-booking-system/utils"
)

// @title Fiber Booking System API
// @version 1.0
// @description A booking system API built with Fiber framework
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@example.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:3000
// @BasePath /api
func main() {
	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Default error handling
			code := fiber.StatusInternalServerError

			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Add global middleware
	app.Use(cors.New())
	app.Use(recover.New())

	// Initialize dependencies
	cache := utils.NewInMemoryCache()
	bookingRepo := repository.NewBookingRepositoryMock()
	bookingUseCase := usecase.NewBookingUseCase(bookingRepo, cache)
	bookingHandler := handler.NewBookingHandler(bookingUseCase)

	// Setup routes
	router.SetupRoutes(app, bookingHandler)

	// Start server
	log.Println("Starting server on :3000")
	log.Fatal(app.Listen("127.0.0.1:3000"))
}
