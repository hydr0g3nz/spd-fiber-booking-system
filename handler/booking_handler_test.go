package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hydr0g3nz/spd-fiber-booking-system/dto"
	"github.com/hydr0g3nz/spd-fiber-booking-system/handler"
	"github.com/hydr0g3nz/spd-fiber-booking-system/mocks"
	"github.com/hydr0g3nz/spd-fiber-booking-system/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupApp(mockUseCase *mocks.BookingUseCase) *fiber.App {
	app := fiber.New()
	bookingHandler := handler.NewBookingHandler(mockUseCase)

	app.Post("/api/bookings", bookingHandler.CreateBooking)
	app.Get("/api/bookings/:id", bookingHandler.GetBooking)
	app.Get("/api/bookings", bookingHandler.GetAllBookings)
	app.Delete("/api/bookings/:id", bookingHandler.CancelBooking)

	return app
}

func TestCreateBookingHandler(t *testing.T) {
	// Create mock use case
	mockUseCase := new(mocks.BookingUseCase)

	// Create test data
	now := time.Now()
	reqPayload := &dto.CreateBookingRequest{
		UserID:    123,
		ServiceID: 456,
		Price:     30000.0,
	}

	createdBooking := &models.Booking{
		ID:        1,
		UserID:    reqPayload.UserID,
		ServiceID: reqPayload.ServiceID,
		Price:     reqPayload.Price,
		Status:    models.BookingStatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Setup expectations
	mockUseCase.On("CreateBooking", mock.Anything, mock.MatchedBy(func(r *dto.CreateBookingRequest) bool {
		return r.UserID == reqPayload.UserID &&
			r.ServiceID == reqPayload.ServiceID &&
			r.Price == reqPayload.Price
	})).Return(createdBooking, nil)

	// Setup app with mock
	app := setupApp(mockUseCase)

	// Create request body
	reqBody, _ := json.Marshal(reqPayload)

	// Perform request
	req := httptest.NewRequest("POST", "/api/bookings", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)

	// Parse response body
	var responseBooking models.Booking
	json.NewDecoder(resp.Body).Decode(&responseBooking)

	assert.Equal(t, int64(1), responseBooking.ID)
	assert.Equal(t, reqPayload.UserID, responseBooking.UserID)
	assert.Equal(t, reqPayload.ServiceID, responseBooking.ServiceID)
	assert.Equal(t, reqPayload.Price, responseBooking.Price)
	assert.Equal(t, models.BookingStatusPending.String(), string(responseBooking.Status))

	mockUseCase.AssertExpectations(t)
}

func TestGetBookingHandler(t *testing.T) {
	// Create mock use case
	mockUseCase := new(mocks.BookingUseCase)

	// Create test data
	bookingID := int64(1)
	booking := &models.Booking{
		ID:        bookingID,
		UserID:    123,
		ServiceID: 456,
		Price:     30000.0,
		Status:    models.BookingStatusConfirmed,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Setup expectations
	mockUseCase.On("GetBookingByID", mock.Anything, bookingID).Return(booking, nil)

	// Setup app with mock
	app := setupApp(mockUseCase)

	// Perform request
	req := httptest.NewRequest("GET", "/api/bookings/1", nil)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Parse response body
	var responseBooking models.Booking
	json.NewDecoder(resp.Body).Decode(&responseBooking)

	assert.Equal(t, booking.ID, responseBooking.ID)
	assert.Equal(t, booking.UserID, responseBooking.UserID)
	assert.Equal(t, booking.ServiceID, responseBooking.ServiceID)
	assert.Equal(t, booking.Price, responseBooking.Price)
	assert.Equal(t, booking.Status.String(), string(responseBooking.Status))

	mockUseCase.AssertExpectations(t)
}

func TestGetAllBookingsHandler(t *testing.T) {
	// Create mock use case
	mockUseCase := new(mocks.BookingUseCase)

	// Create test data
	bookings := []*models.Booking{
		{
			ID:        1,
			UserID:    123,
			ServiceID: 456,
			Price:     30000.0,
			Status:    models.BookingStatusPending,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			UserID:    234,
			ServiceID: 567,
			Price:     45000.0,
			Status:    models.BookingStatusConfirmed,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Setup expectations - match any query parameters
	mockUseCase.On("GetAllBookings", mock.Anything, mock.Anything).Return(bookings, nil)

	// Setup app with mock
	app := setupApp(mockUseCase)

	// Perform request
	req := httptest.NewRequest("GET", "/api/bookings?sort=price&high-value=true", nil)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Parse response body
	var responseBookings []*models.Booking
	json.NewDecoder(resp.Body).Decode(&responseBookings)

	assert.Equal(t, 2, len(responseBookings))
	assert.Equal(t, bookings[0].ID, responseBookings[0].ID)
	assert.Equal(t, bookings[1].ID, responseBookings[1].ID)

	mockUseCase.AssertExpectations(t)
}

func TestCancelBookingHandler_Success(t *testing.T) {
	// Create mock use case
	mockUseCase := new(mocks.BookingUseCase)

	// Create test data
	bookingID := int64(1)
	canceledBooking := &models.Booking{
		ID:        bookingID,
		UserID:    123,
		ServiceID: 456,
		Price:     30000.0,
		Status:    models.BookingStatusCanceled,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Setup expectations
	mockUseCase.On("CancelBooking", mock.Anything, bookingID).Return(canceledBooking, nil)

	// Setup app with mock
	app := setupApp(mockUseCase)

	// Perform request
	req := httptest.NewRequest("DELETE", "/api/bookings/1", nil)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Parse response body
	var responseBooking models.Booking
	json.NewDecoder(resp.Body).Decode(&responseBooking)

	assert.Equal(t, canceledBooking.ID, responseBooking.ID)
	assert.Equal(t, canceledBooking.Status.String(), string(responseBooking.Status))

	mockUseCase.AssertExpectations(t)
}

func TestCancelBookingHandler_CannotCancelConfirmed(t *testing.T) {
	// Create mock use case
	mockUseCase := new(mocks.BookingUseCase)

	// Create test data
	bookingID := int64(1)

	// Setup expectations - confirmed booking cannot be canceled
	mockUseCase.On("CancelBooking", mock.Anything, bookingID).Return(nil,
		fiber.NewError(fiber.StatusBadRequest, "cannot cancel a confirmed booking"))

	// Setup app with mock
	app := setupApp(mockUseCase)

	// Perform request
	req := httptest.NewRequest("DELETE", "/api/bookings/1", nil)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	// Parse response body
	var errorResponse map[string]string
	json.NewDecoder(resp.Body).Decode(&errorResponse)

	assert.Equal(t, "cannot cancel a confirmed booking", errorResponse["error"])

	mockUseCase.AssertExpectations(t)
}

func TestCreateBookingHandler_InvalidRequest(t *testing.T) {
	// Create mock use case
	mockUseCase := new(mocks.BookingUseCase)

	// Invalid request (negative price)
	reqData := map[string]interface{}{
		"user_id":    123,
		"service_id": 456,
		"price":      -100.0, // Invalid negative price
	}

	// Setup app with mock
	app := setupApp(mockUseCase)

	// Create request body
	reqBody, _ := json.Marshal(reqData)

	// Perform request
	req := httptest.NewRequest("POST", "/api/bookings", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode) // Should return bad request

	// Usecase should never be called with invalid data
	mockUseCase.AssertNotCalled(t, "CreateBooking")
}
