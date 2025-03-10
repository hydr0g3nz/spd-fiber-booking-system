package repository

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/hydr0g3nz/spd-fiber-booking-system/models"
)

// BookingRepository defines the interface for booking data operations
type BookingRepository interface {
	Create(ctx context.Context, booking *models.Booking) (*models.Booking, error)
	GetByID(ctx context.Context, id int64) (*models.Booking, error)
	GetAll(ctx context.Context) ([]*models.Booking, error)
	Update(ctx context.Context, booking *models.Booking) (*models.Booking, error)
}

// BookingRepositoryMock is a mock implementation of BookingRepository
type BookingRepositoryMock struct {
	bookings map[int64]*models.Booking
	mutex    sync.RWMutex
	nextID   int64
}

// NewBookingRepositoryMock creates a new instance of BookingRepositoryMock
func NewBookingRepositoryMock() *BookingRepositoryMock {
	repo := &BookingRepositoryMock{
		bookings: make(map[int64]*models.Booking),
		nextID:   11, // Start from 11 since we'll have default bookings 1-10
	}

	// Initialize default bookings (ID 1-10)
	now := time.Now()
	for i := int64(1); i <= 10; i++ {
		status := models.BookingStatusPending
		price := float64(i * 10000) // Prices: 10000, 20000, ... 100000

		// Make some bookings confirmed or rejected for testing
		if i%3 == 0 {
			status = models.BookingStatusConfirmed
		} else if i%5 == 0 {
			status = models.BookingStatusRejected
		}

		booking := &models.Booking{
			ID:        i,
			UserID:    100 + i,
			ServiceID: 200 + i,
			Price:     price,
			Status:    status,
			CreatedAt: now.Add(-time.Duration(i) * time.Hour),
			UpdatedAt: now,
		}
		repo.bookings[i] = booking
	}

	return repo
}

// Create creates a new booking
func (r *BookingRepositoryMock) Create(ctx context.Context, booking *models.Booking) (*models.Booking, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	booking.ID = r.nextID
	r.nextID++
	booking.Status = models.BookingStatusPending

	// Deep copy to avoid reference issues
	newBooking := &models.Booking{
		ID:        booking.ID,
		UserID:    booking.UserID,
		ServiceID: booking.ServiceID,
		Price:     booking.Price,
		Status:    booking.Status,
		CreatedAt: booking.CreatedAt,
		UpdatedAt: booking.UpdatedAt,
	}

	r.bookings[newBooking.ID] = newBooking

	return newBooking, nil
}

// GetByID retrieves a booking by ID
func (r *BookingRepositoryMock) GetByID(ctx context.Context, id int64) (*models.Booking, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	booking, exists := r.bookings[id]
	if !exists {
		return nil, errors.New("booking not found")
	}

	// Return a copy to avoid reference issues
	return &models.Booking{
		ID:        booking.ID,
		UserID:    booking.UserID,
		ServiceID: booking.ServiceID,
		Price:     booking.Price,
		Status:    booking.Status,
		CreatedAt: booking.CreatedAt,
		UpdatedAt: booking.UpdatedAt,
	}, nil
}

// GetAll retrieves all bookings
func (r *BookingRepositoryMock) GetAll(ctx context.Context) ([]*models.Booking, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	bookings := make([]*models.Booking, 0, len(r.bookings))
	for _, booking := range r.bookings {
		// Return copies to avoid reference issues
		bookings = append(bookings, &models.Booking{
			ID:        booking.ID,
			UserID:    booking.UserID,
			ServiceID: booking.ServiceID,
			Price:     booking.Price,
			Status:    booking.Status,
			CreatedAt: booking.CreatedAt,
			UpdatedAt: booking.UpdatedAt,
		})
	}

	return bookings, nil
}

// Update updates a booking
func (r *BookingRepositoryMock) Update(ctx context.Context, booking *models.Booking) (*models.Booking, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	existing, exists := r.bookings[booking.ID]
	if !exists {
		return nil, errors.New("booking not found")
	}

	// Update the booking while preserving creation time
	booking.CreatedAt = existing.CreatedAt

	// Store a copy to avoid reference issues
	updatedBooking := &models.Booking{
		ID:        booking.ID,
		UserID:    booking.UserID,
		ServiceID: booking.ServiceID,
		Price:     booking.Price,
		Status:    booking.Status,
		CreatedAt: booking.CreatedAt,
		UpdatedAt: booking.UpdatedAt,
	}

	r.bookings[booking.ID] = updatedBooking

	return updatedBooking, nil
}
