package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/hydr0g3nz/spd-fiber-booking-system/dto"
	"github.com/hydr0g3nz/spd-fiber-booking-system/usecase"
)

// BookingHandler handles HTTP requests for bookings
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
// @Description Create a new booking with the given details
// @Tags bookings
// @Accept json
// @Produce json
// @Param booking body dto.CreateBookingRequest true "Booking Information"
// @Success 201 {object} models.Booking
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
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
// @Description Get a booking's details by its ID
// @Tags bookings
// @Accept json
// @Produce json
// @Param id path int true "Booking ID"
// @Success 200 {object} models.Booking
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /bookings/{id} [get]
func (h *BookingHandler) GetBooking(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	booking, err := h.bookingUseCase.GetBookingByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Booking not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(booking)
}

// GetAllBookings godoc
// @Summary Get all bookings
// @Description Get all bookings with optional sorting and filtering
// @Tags bookings
// @Accept json
// @Produce json
// @Param sort query string false "Sort by (price or date)"
// @Param high-value query boolean false "Filter high-value bookings (price > 50,000)"
// @Success 200 {array} models.Booking
// @Failure 500 {object} map[string]string
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
// @Description Cancel a booking by its ID
// @Tags bookings
// @Accept json
// @Produce json
// @Param id path int true "Booking ID"
// @Success 200 {object} models.Booking
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /bookings/{id} [delete]
func (h *BookingHandler) CancelBooking(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	booking, err := h.bookingUseCase.CancelBooking(c.Context(), id)
	if err != nil {
		if err.Error() == "cannot cancel a confirmed booking" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Booking not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(booking)
}
