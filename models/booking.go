package models

import (
	"time"
)

// Booking represents a booking entity
// @Description Booking entity representing a customer's service booking
type Booking struct {
	ID        int64         `json:"id" example:"1" description:"Booking ID"`
	UserID    int64         `json:"user_id" example:"123" description:"User ID"`
	ServiceID int64         `json:"service_id" example:"456" description:"Service ID"`
	Price     float64       `json:"price" example:"30000.0" description:"Booking price"`
	Status    BookingStatus `json:"status" enums:"pending,confirmed,rejected,canceled" example:"pending" description:"Booking status"`
	CreatedAt time.Time     `json:"created_at" format:"date-time" example:"2024-03-11T12:00:00Z" description:"Creation timestamp"`
	UpdatedAt time.Time     `json:"updated_at" format:"date-time" example:"2024-03-11T12:00:00Z" description:"Last update timestamp"`
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
