package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hydr0g3nz/spd-fiber-booking-system/dto"
	"github.com/hydr0g3nz/spd-fiber-booking-system/usecase"
)

// SwaggerModels defines models for swagger documentation
// @title Fiber Booking System API
// @version 1.0
// @description A booking management system built with Go and Fiber framework
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@example.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:3000
// @BasePath /api
// @schemes http

// BookingHandler manages HTTP requests for booking endpoints
type BookingHandler struct {
	bookingUseCase usecase.BookingUseCase
}

// NewBookingHandler creates a new instance of BookingHandler
func NewBookingHandler(bookingUseCase usecase.BookingUseCase) *BookingHandler {
	return &BookingHandler{
		bookingUseCase: bookingUseCase,
	}
}

// CreateBooking godoc
// @Summary Create a new booking
// @Description Create a new booking with the provided details
// @Tags bookings
// @Accept json
// @Produce json
// @Param booking body dto.CreateBookingRequest true "Booking Information"
// @Success 201 {object} models.Booking "Created booking"
// @Failure 400 {object} map[string]string "Invalid request parameters"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /bookings [post]
func (h *BookingHandler) CreateBooking(c *fiber.Ctx) error {
	req := new(dto.CreateBookingRequest)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if req.UserID <= 0 || req.ServiceID <= 0 || req.Price <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "UserID, ServiceID, and Price are required and must be positive values",
		})
	}

	booking, err := h.bookingUseCase.CreateBooking(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(booking)
}

// GetBooking godoc
// @Summary Get a booking by ID
// @Description Get detailed information about a specific booking
// @Tags bookings
// @Accept json
// @Produce json
// @Param id path int true "Booking ID" minimum(1)
// @Success 200 {object} models.Booking "Booking details"
// @Failure 400 {object} map[string]string "Invalid booking ID format"
// @Failure 404 {object} map[string]string "Booking not found"
// @Router /bookings/{id} [get]
func (h *BookingHandler) GetBooking(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid booking ID format",
		})
	}

	booking, err := h.bookingUseCase.GetBookingByID(c.Context(), int64(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Booking not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(booking)
}

// GetAllBookings godoc
// @Summary Get all bookings
// @Description Get a list of all bookings with optional sorting and filtering
// @Tags bookings
// @Accept json
// @Produce json
// @Param sort query string false "Sort by field (price or date)"
// @Param high-value query boolean false "Filter high-value bookings (price > 50,000)"
// @Success 200 {array} models.Booking "List of bookings"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /bookings [get]
func (h *BookingHandler) GetAllBookings(c *fiber.Ctx) error {
	params := &dto.BookingsQueryParams{
		Sort:      c.Query("sort"),
		HighValue: c.Query("high-value") == "true",
	}

	bookings, err := h.bookingUseCase.GetAllBookings(c.Context(), params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(bookings)
}

// CancelBooking godoc
// @Summary Cancel a booking
// @Description Cancel an existing booking by its ID
// @Tags bookings
// @Accept json
// @Produce json
// @Param id path int true "Booking ID" minimum(1)
// @Success 200 {object} models.Booking "Canceled booking details"
// @Failure 400 {object} map[string]string "Invalid booking ID or cannot cancel"
// @Failure 404 {object} map[string]string "Booking not found"
// @Router /bookings/{id} [delete]
func (h *BookingHandler) CancelBooking(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid booking ID format",
		})
	}

	booking, err := h.bookingUseCase.CancelBooking(c.Context(), int64(id))
	if err != nil {
		// Check if this is a business rule error (cannot cancel confirmed booking)
		if err.Error() == "cannot cancel a confirmed booking" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		// Otherwise, it's likely a "not found" error
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Booking not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(booking)
}
