package utils_test

import (
	"testing"
	"time"

	"github.com/hydr0g3nz/spd-fiber-booking-system/models"
	"github.com/hydr0g3nz/spd-fiber-booking-system/utils"
	"github.com/stretchr/testify/assert"
)

func TestSortBookingsByPriceAscending(t *testing.T) {
	// Arrange - Create sample bookings with different prices
	bookings := []*models.Booking{
		{ID: 1, Price: 30000.0},
		{ID: 2, Price: 10000.0},
		{ID: 3, Price: 50000.0},
		{ID: 4, Price: 20000.0},
	}

	// Act - Sort by price in ascending order
	utils.SortBookingsByPrice(bookings, true)

	// Assert - Check if sorted correctly
	assert.Equal(t, int64(2), bookings[0].ID, "First booking should be ID 2 (lowest price)")
	assert.Equal(t, int64(4), bookings[1].ID, "Second booking should be ID 4")
	assert.Equal(t, int64(1), bookings[2].ID, "Third booking should be ID 1")
	assert.Equal(t, int64(3), bookings[3].ID, "Fourth booking should be ID 3 (highest price)")

	// Verify prices are in ascending order
	for i := 0; i < len(bookings)-1; i++ {
		assert.LessOrEqual(t, bookings[i].Price, bookings[i+1].Price,
			"Prices should be in ascending order")
	}
}

func TestSortBookingsByPriceDescending(t *testing.T) {
	// Arrange - Create sample bookings with different prices
	bookings := []*models.Booking{
		{ID: 1, Price: 30000.0},
		{ID: 2, Price: 10000.0},
		{ID: 3, Price: 50000.0},
		{ID: 4, Price: 20000.0},
	}

	// Act - Sort by price in descending order
	utils.SortBookingsByPrice(bookings, false)

	// Assert - Check if sorted correctly
	assert.Equal(t, int64(3), bookings[0].ID, "First booking should be ID 3 (highest price)")
	assert.Equal(t, int64(1), bookings[1].ID, "Second booking should be ID 1")
	assert.Equal(t, int64(4), bookings[2].ID, "Third booking should be ID 4")
	assert.Equal(t, int64(2), bookings[3].ID, "Fourth booking should be ID 2 (lowest price)")

	// Verify prices are in descending order
	for i := 0; i < len(bookings)-1; i++ {
		assert.GreaterOrEqual(t, bookings[i].Price, bookings[i+1].Price,
			"Prices should be in descending order")
	}
}

func TestSortBookingsByDateAscending(t *testing.T) {
	// Arrange - Create sample bookings with different dates
	now := time.Now()
	bookings := []*models.Booking{
		{ID: 1, CreatedAt: now.Add(-2 * time.Hour)},  // 2 hours ago
		{ID: 2, CreatedAt: now.Add(-24 * time.Hour)}, // 1 day ago
		{ID: 3, CreatedAt: now},                      // now
		{ID: 4, CreatedAt: now.Add(-12 * time.Hour)}, // 12 hours ago
	}

	// Act - Sort by date in ascending order (oldest first)
	utils.SortBookingsByDate(bookings, true)

	// Assert - Check if sorted correctly
	assert.Equal(t, int64(2), bookings[0].ID, "First booking should be ID 2 (oldest)")
	assert.Equal(t, int64(4), bookings[1].ID, "Second booking should be ID 4")
	assert.Equal(t, int64(1), bookings[2].ID, "Third booking should be ID 1")
	assert.Equal(t, int64(3), bookings[3].ID, "Fourth booking should be ID 3 (newest)")

	// Verify dates are in ascending order
	for i := 0; i < len(bookings)-1; i++ {
		assert.True(t, bookings[i].CreatedAt.Before(bookings[i+1].CreatedAt) ||
			bookings[i].CreatedAt.Equal(bookings[i+1].CreatedAt),
			"Dates should be in ascending order")
	}
}

func TestSortBookingsByDateDescending(t *testing.T) {
	// Arrange - Create sample bookings with different dates
	now := time.Now()
	bookings := []*models.Booking{
		{ID: 1, CreatedAt: now.Add(-2 * time.Hour)},  // 2 hours ago
		{ID: 2, CreatedAt: now.Add(-24 * time.Hour)}, // 1 day ago
		{ID: 3, CreatedAt: now},                      // now
		{ID: 4, CreatedAt: now.Add(-12 * time.Hour)}, // 12 hours ago
	}

	// Act - Sort by date in descending order (newest first)
	utils.SortBookingsByDate(bookings, false)

	// Assert - Check if sorted correctly
	assert.Equal(t, int64(3), bookings[0].ID, "First booking should be ID 3 (newest)")
	assert.Equal(t, int64(1), bookings[1].ID, "Second booking should be ID 1")
	assert.Equal(t, int64(4), bookings[2].ID, "Third booking should be ID 4")
	assert.Equal(t, int64(2), bookings[3].ID, "Fourth booking should be ID 2 (oldest)")

	// Verify dates are in descending order
	for i := 0; i < len(bookings)-1; i++ {
		assert.True(t, bookings[i].CreatedAt.After(bookings[i+1].CreatedAt) ||
			bookings[i].CreatedAt.Equal(bookings[i+1].CreatedAt),
			"Dates should be in descending order")
	}
}

func TestFilterHighValueBookings(t *testing.T) {
	// Arrange - Create sample bookings with different prices
	bookings := []*models.Booking{
		{ID: 1, Price: 30000.0},  // Not high value
		{ID: 2, Price: 60000.0},  // High value
		{ID: 3, Price: 50000.1},  // High value (just over threshold)
		{ID: 4, Price: 50000.0},  // Not high value (exactly at threshold)
		{ID: 5, Price: 100000.0}, // High value
	}

	// Act - Filter high value bookings
	highValueBookings := utils.FilterHighValueBookings(bookings)

	// Assert - Check if filtered correctly
	assert.Equal(t, 3, len(highValueBookings), "Should have 3 high value bookings")

	// Verify all returned bookings are high value
	for _, booking := range highValueBookings {
		assert.Greater(t, booking.Price, 50000.0, "All filtered bookings should have price > 50000")
	}

	// Verify the IDs of the high value bookings
	highValueIDs := map[int64]bool{}
	for _, b := range highValueBookings {
		highValueIDs[b.ID] = true
	}

	assert.True(t, highValueIDs[2], "Booking with ID 2 should be included")
	assert.True(t, highValueIDs[3], "Booking with ID 3 should be included")
	assert.True(t, highValueIDs[5], "Booking with ID 5 should be included")
	assert.False(t, highValueIDs[1], "Booking with ID 1 should not be included")
	assert.False(t, highValueIDs[4], "Booking with ID 4 should not be included")
}

func TestFilterHighValueBookings_EmptyInput(t *testing.T) {
	// Arrange - Empty bookings slice
	var bookings []*models.Booking

	// Act - Filter high value bookings
	highValueBookings := utils.FilterHighValueBookings(bookings)

	// Assert - Should return empty slice, not nil
	assert.NotNil(t, highValueBookings, "Should return empty slice, not nil")
	assert.Equal(t, 0, len(highValueBookings), "Should have 0 bookings")
}

func TestFilterHighValueBookings_NoHighValueBookings(t *testing.T) {
	// Arrange - Create sample bookings with no high value ones
	bookings := []*models.Booking{
		{ID: 1, Price: 30000.0},
		{ID: 2, Price: 10000.0},
		{ID: 3, Price: 50000.0}, // Exactly at threshold
		{ID: 4, Price: 49999.9}, // Just under threshold
	}

	// Act - Filter high value bookings
	highValueBookings := utils.FilterHighValueBookings(bookings)

	// Assert - Should return empty slice
	assert.Equal(t, 0, len(highValueBookings), "Should have 0 high value bookings")
}

func TestFilterHighValueBookings_AllHighValueBookings(t *testing.T) {
	// Arrange - Create sample bookings that are all high value
	bookings := []*models.Booking{
		{ID: 1, Price: 50000.1},
		{ID: 2, Price: 60000.0},
		{ID: 3, Price: 100000.0},
	}

	// Act - Filter high value bookings
	highValueBookings := utils.FilterHighValueBookings(bookings)

	// Assert - Should return all bookings
	assert.Equal(t, len(bookings), len(highValueBookings), "Should return all bookings")

	// Verify all bookings are present
	assert.Equal(t, bookings, highValueBookings, "Should contain same bookings in same order")
}
