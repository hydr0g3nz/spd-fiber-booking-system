package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/hydr0g3nz/spd-fiber-booking-system/dto"
	"github.com/hydr0g3nz/spd-fiber-booking-system/models"
	"github.com/hydr0g3nz/spd-fiber-booking-system/repository"
	"github.com/hydr0g3nz/spd-fiber-booking-system/utils"
)

// BookingUseCase defines the interface for booking business logic
type BookingUseCase interface {
	CreateBooking(ctx context.Context, req *dto.CreateBookingRequest) (*models.Booking, error)
	GetBookingByID(ctx context.Context, id int64) (*models.Booking, error)
	GetAllBookings(ctx context.Context, params *dto.BookingsQueryParams) ([]*models.Booking, error)
	CancelBooking(ctx context.Context, id int64) (*models.Booking, error)
}

// BookingUseCaseImpl implements BookingUseCase
type BookingUseCaseImpl struct {
	repo   repository.BookingRepository
	cache  utils.Cache
	nextID int64
	mu     sync.Mutex
}

// NewBookingUseCase creates a new instance of BookingUseCaseImpl
func NewBookingUseCase(repo repository.BookingRepository, cache utils.Cache) BookingUseCase {
	uc := &BookingUseCaseImpl{
		repo:   repo,
		cache:  cache,
		nextID: 11, // Start from 11 since we'll have default bookings 1-10
	}

	// Start background task to check for expired bookings
	go uc.checkExpiredBookings()

	return uc
}

// CreateBooking creates a new booking
func (uc *BookingUseCaseImpl) CreateBooking(ctx context.Context, req *dto.CreateBookingRequest) (*models.Booking, error) {
	uc.mu.Lock() // ล็อคก่อนอ่าน/เขียนค่า nextID
	id := uc.nextID
	uc.nextID++
	uc.mu.Unlock()
	booking := &models.Booking{
		ID:        id,
		UserID:    req.UserID,
		ServiceID: req.ServiceID,
		Price:     req.Price,
		Status:    models.BookingStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// // Save booking to repository
	// newBooking, err := uc.repo.Create(ctx, booking)
	// if err != nil {
	// 	return nil, err
	// }

	// Store in cache
	cacheKey := fmt.Sprintf("booking:%d", booking.ID)
	uc.cache.Set(cacheKey, booking)

	// For high-value bookings, run credit check in background
	if booking.Price > 50000 {
		go uc.checkCredit(booking)
	}

	return booking, nil
}

// GetBookingByID retrieves a booking by ID
func (uc *BookingUseCaseImpl) GetBookingByID(ctx context.Context, id int64) (*models.Booking, error) {
	// Try to get from cache first
	cacheKey := fmt.Sprintf("booking:%d", id)
	if cachedValue, found := uc.cache.Get(cacheKey); found {
		fmt.Println("Get from cache")
		return cachedValue.(*models.Booking), nil
	}

	// If not in cache, get from repository
	booking, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	fmt.Println("Get from repository")
	return booking, nil
}

// GetAllBookings retrieves all bookings with optional sorting and filtering
func (uc *BookingUseCaseImpl) GetAllBookings(ctx context.Context, params *dto.BookingsQueryParams) ([]*models.Booking, error) {
	// Get all bookings from repository
	bookingsRepo, err := uc.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	bookingsCache := uc.cache.GetAll()
	mergedBookings := make([]*models.Booking, 0, len(bookingsRepo)+len(bookingsCache))
	for _, booking := range bookingsCache {
		mergedBookings = append(mergedBookings, booking.(*models.Booking))
	}
	mergedBookings = append(mergedBookings, bookingsRepo...)
	// Merge repository and cache bookings
	// Filter high-value bookings if requested
	if params.HighValue {
		mergedBookings = utils.FilterHighValueBookings(mergedBookings)
	}

	// Sort bookings if requested
	switch params.Sort {
	case "price":
		utils.SortBookingsByPrice(mergedBookings, true)
	case "date":
		utils.SortBookingsByDate(mergedBookings, true)
	}

	return mergedBookings, nil
}

// CancelBooking cancels a booking
func (uc *BookingUseCaseImpl) CancelBooking(ctx context.Context, id int64) (*models.Booking, error) {
	// Get booking from repository
	booking, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cannot cancel a confirmed booking
	if booking.Status == models.BookingStatusConfirmed {
		return nil, errors.New("cannot cancel a confirmed booking")
	}

	// Update status to canceled
	booking.Status = models.BookingStatusCanceled
	booking.UpdatedAt = time.Now()

	// Update in repository
	updatedBooking, err := uc.repo.Update(ctx, booking)
	if err != nil {
		return nil, err
	}

	// Remove from cache
	cacheKey := fmt.Sprintf("booking:%d", id)
	uc.cache.Delete(cacheKey)

	return updatedBooking, nil
}

// checkCredit simulates a credit check for high-value bookings
func (uc *BookingUseCaseImpl) checkCredit(booking *models.Booking) {

	// Random credit check result (70% success rate)
	rand.Seed(time.Now().UnixNano())
	status := models.BookingStatusConfirmed
	if rand.Float64() < 0.3 { // 30% chance of rejection
		status = models.BookingStatusRejected
	}

	// Update booking status
	booking.Status = status
	booking.UpdatedAt = time.Now()

	cacheKey := fmt.Sprintf("booking:%d", booking.ID)
	uc.cache.Set(cacheKey, booking)

	log.Printf("Credit check completed for booking %d. Status: %s", booking.ID, status)
}

// checkExpiredBookings runs as a background task to cancel expired bookings
func (uc *BookingUseCaseImpl) checkExpiredBookings() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("Running expired bookings check...")

		// Get all bookings from repository
		ctx := context.Background()
		bookings, err := uc.repo.GetAll(ctx)
		if err != nil {
			log.Printf("Error fetching bookings: %v", err)
			continue
		}

		now := time.Now()
		expiredCount := 0

		for _, booking := range bookings {
			// If booking is pending for more than 5 minutes, mark as canceled
			if booking.Status == models.BookingStatusPending && now.Sub(booking.CreatedAt) > 5*time.Minute {
				booking.Status = models.BookingStatusCanceled
				booking.UpdatedAt = now

				// Update in repository
				_, err := uc.repo.Update(ctx, booking)
				if err != nil {
					log.Printf("Error updating expired booking %d: %v", booking.ID, err)
					continue
				}

				expiredCount++
			}
		}
		// Get all bookings from cache
		cacheBookings := uc.cache.GetAll()
		for cacheKey, cachedValue := range cacheBookings {
			booking := cachedValue.(*models.Booking)

			// If booking is pending for more than 5 minutes, mark as canceled
			if booking.Status == models.BookingStatusPending && now.Sub(booking.CreatedAt) > 5*time.Minute {
				booking.Status = models.BookingStatusCanceled
				booking.UpdatedAt = now

				// Update in cache
				uc.cache.Set(cacheKey, booking)

				expiredCount++
			}
		}

		if expiredCount > 0 {
			log.Printf("Auto-canceled %d expired bookings", expiredCount)
		}
	}
}
