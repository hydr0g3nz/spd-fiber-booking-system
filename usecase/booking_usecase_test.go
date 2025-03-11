package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/hydr0g3nz/spd-fiber-booking-system/dto"
	"github.com/hydr0g3nz/spd-fiber-booking-system/mocks"
	"github.com/hydr0g3nz/spd-fiber-booking-system/models"
	"github.com/hydr0g3nz/spd-fiber-booking-system/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateBooking(t *testing.T) {
	// Create mocks
	mockRepo := new(mocks.BookingRepository)
	mockCache := new(mocks.Cache)

	// Create test data
	now := time.Now()
	req := &dto.CreateBookingRequest{
		UserID:    123,
		ServiceID: 456,
		Price:     30000.0,
	}

	// expectedBooking := &models.Booking{
	// 	UserID:    req.UserID,
	// 	ServiceID: req.ServiceID,
	// 	Price:     req.Price,
	// 	Status:    models.BookingStatusPending,
	// 	CreatedAt: now,
	// 	UpdatedAt: now,
	// }

	createdBooking := &models.Booking{
		ID:        1,
		UserID:    req.UserID,
		ServiceID: req.ServiceID,
		Price:     req.Price,
		Status:    models.BookingStatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Setup expectations
	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(b *models.Booking) bool {
		return b.UserID == req.UserID &&
			b.ServiceID == req.ServiceID &&
			b.Price == req.Price &&
			b.Status == models.BookingStatusPending
	})).Return(createdBooking, nil)

	mockCache.On("Set", "booking:1", createdBooking).Return()

	// Create use case
	uc := usecase.NewBookingUseCase(mockRepo, mockCache)

	// Execute
	result, err := uc.CreateBooking(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, req.UserID, result.UserID)
	assert.Equal(t, req.ServiceID, result.ServiceID)
	assert.Equal(t, req.Price, result.Price)
	assert.Equal(t, models.BookingStatusPending, result.Status)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestGetBookingByID_FromCache(t *testing.T) {
	// Create mocks
	mockRepo := new(mocks.BookingRepository)
	mockCache := new(mocks.Cache)

	// Create test data
	bookingID := int64(1)
	booking := &models.Booking{
		ID:        bookingID,
		UserID:    123,
		ServiceID: 456,
		Price:     30000.0,
		Status:    models.BookingStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Setup expectations - booking found in cache
	mockCache.On("Get", "booking:1").Return(booking, true)

	// Create use case
	uc := usecase.NewBookingUseCase(mockRepo, mockCache)

	// Execute
	result, err := uc.GetBookingByID(context.Background(), bookingID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, booking, result)

	// Verify that repository was not called (cache hit)
	mockRepo.AssertNotCalled(t, "GetByID")
	mockCache.AssertExpectations(t)
}

func TestGetBookingByID_FromRepository(t *testing.T) {
	// Create mocks
	mockRepo := new(mocks.BookingRepository)
	mockCache := new(mocks.Cache)

	// Create test data
	bookingID := int64(1)
	booking := &models.Booking{
		ID:        bookingID,
		UserID:    123,
		ServiceID: 456,
		Price:     30000.0,
		Status:    models.BookingStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Setup expectations - booking not found in cache, but found in repository
	mockCache.On("Get", "booking:1").Return(nil, false)
	mockRepo.On("GetByID", mock.Anything, bookingID).Return(booking, nil)
	mockCache.On("Set", "booking:1", booking).Return()

	// Create use case
	uc := usecase.NewBookingUseCase(mockRepo, mockCache)

	// Execute
	result, err := uc.GetBookingByID(context.Background(), bookingID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, booking, result)

	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestGetAllBookings(t *testing.T) {
	// Create mocks
	mockRepo := new(mocks.BookingRepository)
	mockCache := new(mocks.Cache)

	// Create test data
	params := &dto.BookingsQueryParams{
		Sort:      "price",
		HighValue: false,
	}

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

	cacheMap := map[string]interface{}{
		"booking:1": bookings[0],
		"booking:2": bookings[1],
	}

	// Setup expectations
	mockRepo.On("GetAll", mock.Anything).Return(bookings, nil)
	mockCache.On("GetAll").Return(cacheMap)

	// Create use case
	uc := usecase.NewBookingUseCase(mockRepo, mockCache)

	// Execute
	result, err := uc.GetAllBookings(context.Background(), params)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	// Since we specified sort by price, the first item should be the one with lower price
	assert.Equal(t, float64(30000.0), result[0].Price)
	assert.Equal(t, float64(45000.0), result[1].Price)

	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestCancelBooking_Success(t *testing.T) {
	// Create mocks
	mockRepo := new(mocks.BookingRepository)
	mockCache := new(mocks.Cache)

	// Create test data
	bookingID := int64(1)
	booking := &models.Booking{
		ID:        bookingID,
		UserID:    123,
		ServiceID: 456,
		Price:     30000.0,
		Status:    models.BookingStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	canceledBooking := &models.Booking{
		ID:        bookingID,
		UserID:    123,
		ServiceID: 456,
		Price:     30000.0,
		Status:    models.BookingStatusCanceled,
		CreatedAt: booking.CreatedAt,
		UpdatedAt: time.Now(),
	}

	// Setup expectations
	mockCache.On("Get", "booking:1").Return(booking, true)
	mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(b *models.Booking) bool {
		return b.ID == bookingID && b.Status == models.BookingStatusCanceled
	})).Return(canceledBooking, nil)
	mockCache.On("Set", "booking:1", canceledBooking).Return()

	// Create use case
	uc := usecase.NewBookingUseCase(mockRepo, mockCache)

	// Execute
	result, err := uc.CancelBooking(context.Background(), bookingID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, models.BookingStatusCanceled, result.Status)

	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestCancelBooking_CannotCancelConfirmed(t *testing.T) {
	// Create mocks
	mockRepo := new(mocks.BookingRepository)
	mockCache := new(mocks.Cache)

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
	mockCache.On("Get", "booking:1").Return(booking, true)

	// Create use case
	uc := usecase.NewBookingUseCase(mockRepo, mockCache)

	// Execute
	result, err := uc.CancelBooking(context.Background(), bookingID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "cannot cancel a confirmed booking", err.Error())

	// The repository update should not be called since cancellation fails
	mockRepo.AssertNotCalled(t, "Update")
	mockCache.AssertExpectations(t)
}
