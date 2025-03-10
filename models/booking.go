package models

import (
	"time"
)

// Booking represents a booking entity
type Booking struct {
	ID        int64         `json:"id"`
	UserID    int64         `json:"user_id"`
	ServiceID int64         `json:"service_id"`
	Price     float64       `json:"price"`
	Status    BookingStatus `json:"status"` // pending, confirmed, rejected, canceled
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

// BookingStatus represents the status of a booking as a string type
type BookingStatus string

// BookingStatus constants
const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusRejected  BookingStatus = "rejected"
	BookingStatusCanceled  BookingStatus = "canceled"
)

// IsValid checks if the booking status is valid
func (s BookingStatus) IsValid() bool {
	switch s {
	case BookingStatusPending, BookingStatusConfirmed, BookingStatusRejected, BookingStatusCanceled:
		return true
	}
	return false
}

// String returns the string representation of the booking status
func (s BookingStatus) String() string {
	return string(s)
}
