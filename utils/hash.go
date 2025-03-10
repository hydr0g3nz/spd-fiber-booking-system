package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strconv"

	"github.com/hydr0g3nz/spd-fiber-booking-system/models"
)

// GenerateBookingHash creates a hash for a booking based on its key fields
func GenerateBookingHash(userID, serviceID int64, price float64) string {
	data := strconv.FormatInt(userID, 10) +
		strconv.FormatInt(serviceID, 10) +
		strconv.FormatFloat(price, 'f', 2, 64)

	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// SortBookingsByPrice sorts a slice of booking pointers by price
func SortBookingsByPrice(bookings []*models.Booking, ascending bool) {
	sort.Slice(bookings, func(i, j int) bool {
		if ascending {
			return bookings[i].Price < bookings[j].Price
		}
		return bookings[i].Price > bookings[j].Price
	})
}

// SortBookingsByDate sorts a slice of booking pointers by creation date
func SortBookingsByDate(bookings []*models.Booking, ascending bool) {
	sort.Slice(bookings, func(i, j int) bool {
		if ascending {
			return bookings[i].CreatedAt.Before(bookings[j].CreatedAt)
		}
		return bookings[i].CreatedAt.After(bookings[j].CreatedAt)
	})
}

// FilterHighValueBookings filters bookings with price > 50000
func FilterHighValueBookings(bookings []*models.Booking) []*models.Booking {
	const highValueThreshold = 50000.0

	result := make([]*models.Booking, 0)
	for _, booking := range bookings {
		if booking.Price > highValueThreshold {
			result = append(result, booking)
		}
	}

	return result
}
