package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/hydr0g3nz/spd-fiber-booking-system/handler"
	"github.com/hydr0g3nz/spd-fiber-booking-system/middleware"
)

// SetupRoutes configures all application routes
func SetupRoutes(app *fiber.App, bookingHandler *handler.BookingHandler) {
	// Swagger documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	// API group
	api := app.Group("/api")

	// Apply global middleware
	api.Use(middleware.Logging())

	// Bookings endpoints
	bookings := api.Group("/bookings")
	bookings.Post("/", bookingHandler.CreateBooking)
	bookings.Get("/", bookingHandler.GetAllBookings)
	bookings.Get("/:id", bookingHandler.GetBooking)
	bookings.Delete("/:id", bookingHandler.CancelBooking)

	// Root route for API - redirect to Swagger docs
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/swagger/index.html")
	})
}
