package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/hydr0g3nz/spd-fiber-booking-system/models"
	"github.com/hydr0g3nz/spd-fiber-booking-system/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BookingRepositoryTestSuite struct {
	suite.Suite
	repo repository.BookingRepository
}

func (suite *BookingRepositoryTestSuite) SetupTest() {
	// Create a new repository instance for each test
	suite.repo = repository.NewBookingRepositoryMock()
}

func (suite *BookingRepositoryTestSuite) TestCreate() {
	// Create test data
	ctx := context.Background()
	now := time.Now()
	booking := &models.Booking{
		UserID:    999,
		ServiceID: 888,
		Price:     25000.0,
		Status:    models.BookingStatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Execute
	result, err := suite.repo.Create(ctx, booking)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Greater(suite.T(), result.ID, int64(0))
	assert.Equal(suite.T(), booking.UserID, result.UserID)
	assert.Equal(suite.T(), booking.ServiceID, result.ServiceID)
	assert.Equal(suite.T(), booking.Price, result.Price)
	assert.Equal(suite.T(), models.BookingStatusPending, result.Status)
}

func (suite *BookingRepositoryTestSuite) TestGetByID_ExistingBooking() {
	// Execute - try to get a booking we know exists (ID 1)
	ctx := context.Background()
	result, err := suite.repo.GetByID(ctx, 1)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), int64(1), result.ID)
}

func (suite *BookingRepositoryTestSuite) TestGetByID_NonExistingBooking() {
	// Execute - try to get a booking that doesn't exist
	ctx := context.Background()
	result, err := suite.repo.GetByID(ctx, 999)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), "booking not found", err.Error())
}

func (suite *BookingRepositoryTestSuite) TestGetAll() {
	// Execute
	ctx := context.Background()
	results, err := suite.repo.GetAll(ctx)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), results)
	assert.GreaterOrEqual(suite.T(), len(results), 10) // We should have at least 10 default bookings
}

func (suite *BookingRepositoryTestSuite) TestUpdate_ExistingBooking() {
	// Setup
	ctx := context.Background()
	existingBooking, _ := suite.repo.GetByID(ctx, 1)

	// Modify booking
	updatedBooking := &models.Booking{
		ID:        existingBooking.ID,
		UserID:    existingBooking.UserID,
		ServiceID: existingBooking.ServiceID,
		Price:     existingBooking.Price,
		Status:    models.BookingStatusCanceled, // Change status
		CreatedAt: existingBooking.CreatedAt,
		UpdatedAt: time.Now(),
	}

	// Execute
	result, err := suite.repo.Update(ctx, updatedBooking)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), models.BookingStatusCanceled, result.Status)

	// Verify update was persisted
	retrievedBooking, _ := suite.repo.GetByID(ctx, 1)
	assert.Equal(suite.T(), models.BookingStatusCanceled, retrievedBooking.Status)
}

func (suite *BookingRepositoryTestSuite) TestUpdate_NonExistingBooking() {
	// Setup
	ctx := context.Background()
	nonExistingBooking := &models.Booking{
		ID:        999, // Non-existing ID
		UserID:    123,
		ServiceID: 456,
		Price:     30000.0,
		Status:    models.BookingStatusCanceled,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Execute
	result, err := suite.repo.Update(ctx, nonExistingBooking)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), "booking not found", err.Error())
}

// Run the test suite
func TestBookingRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(BookingRepositoryTestSuite))
}
