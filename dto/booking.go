package dto

import "time"

// CreateBookingRequest is the DTO for creating a new booking
type CreateBookingRequest struct {
	UserID    int64   `json:"user_id" validate:"required"`
	ServiceID int64   `json:"service_id" validate:"required"`
	Price     float64 `json:"price" validate:"required"`
}

// BookingResponse is the DTO for returning booking information
type BookingResponse struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	ServiceID int64     `json:"service_id"`
	Price     float64   `json:"price"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BookingsQueryParams represents query parameters for listing bookings
type BookingsQueryParams struct {
	Sort      string `query:"sort"`
	HighValue bool   `query:"high-value"`
}
