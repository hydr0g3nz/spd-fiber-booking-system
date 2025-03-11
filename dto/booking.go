package dto

import "time"

// CreateBookingRequest is the DTO for creating a new booking
// @Description Request payload for creating a new booking
type CreateBookingRequest struct {
	UserID    int64   `json:"user_id" validate:"required" example:"123" description:"User ID"`
	ServiceID int64   `json:"service_id" validate:"required" example:"456" description:"Service ID"`
	Price     float64 `json:"price" validate:"required" example:"30000.0" description:"Booking price"`
}

// BookingResponse is the DTO for returning booking information
// @Description Response payload for booking information
type BookingResponse struct {
	ID        int64     `json:"id" example:"1" description:"Booking ID"`
	UserID    int64     `json:"user_id" example:"123" description:"User ID"`
	ServiceID int64     `json:"service_id" example:"456" description:"Service ID"`
	Price     float64   `json:"price" example:"30000.0" description:"Booking price"`
	Status    string    `json:"status" example:"pending" description:"Booking status (pending, confirmed, rejected, canceled)"`
	CreatedAt time.Time `json:"created_at" format:"date-time" example:"2024-03-11T12:00:00Z" description:"Creation timestamp"`
	UpdatedAt time.Time `json:"updated_at" format:"date-time" example:"2024-03-11T12:00:00Z" description:"Last update timestamp"`
}

// BookingsQueryParams represents query parameters for listing bookings
// @Description Query parameters for filtering and sorting bookings
type BookingsQueryParams struct {
	Sort      string `query:"sort" example:"price" description:"Sort by field (price or date)"`
	HighValue bool   `query:"high-value" example:"true" description:"Filter high-value bookings (price > 50,000)"`
}
