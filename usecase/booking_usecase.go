package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
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
	repo  repository.BookingRepository
	cache utils.Cache
}

// NewBookingUseCase creates a new instance of BookingUseCaseImpl
func NewBookingUseCase(repo repository.BookingRepository, cache utils.Cache) BookingUseCase {
	uc := &BookingUseCaseImpl{
		repo:  repo,
		cache: cache,
	}

	// Start background task to check for expired bookings
	go uc.checkExpiredBookings()

	return uc
}

// CreateBooking creates a new booking
func (uc *BookingUseCaseImpl) CreateBooking(ctx context.Context, req *dto.CreateBookingRequest) (*models.Booking, error) {
	booking := &models.Booking{
		UserID:    req.UserID,
		ServiceID: req.ServiceID,
		Price:     req.Price,
		Status:    models.BookingStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save booking to repository
	newBooking, err := uc.repo.Create(ctx, booking)
	if err != nil {
		return nil, err
	}

	// Store in cache
	cacheKey := fmt.Sprintf("booking:%d", newBooking.ID)
	uc.cache.Set(cacheKey, newBooking)

	// For high-value bookings, run credit check in background
	if newBooking.Price > 50000 {
		go uc.checkCredit(ctx, newBooking)
	}

	return newBooking, nil
}

// GetBookingByID retrieves a booking by ID
func (uc *BookingUseCaseImpl) GetBookingByID(ctx context.Context, id int64) (*models.Booking, error) {
	// Try to get from cache first
	cacheKey := fmt.Sprintf("booking:%d", id)
	if cachedValue, found := uc.cache.Get(cacheKey); found {
		log.Println("Booking retrieved from cache")
		return cachedValue.(*models.Booking), nil
	}

	// If not in cache, get from repository
	booking, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Add to cache for future use
	uc.cache.Set(cacheKey, booking)
	log.Println("Booking retrieved from repository and added to cache")

	return booking, nil
}

// GetAllBookings retrieves all bookings with optional sorting and filtering
func (uc *BookingUseCaseImpl) GetAllBookings(ctx context.Context, params *dto.BookingsQueryParams) ([]*models.Booking, error) {
	// Get all bookings from repository - this is our source of truth
	bookings, err := uc.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// Check cache for any new bookings not yet in repository
	cacheBookings := uc.cache.GetAll()
	bookingMap := make(map[int64]*models.Booking)

	// First, add all bookings from repository to our map
	for _, booking := range bookings {
		bookingMap[booking.ID] = booking
	}

	// Then, check for any bookings in cache that might not be in repository yet
	for _, cachedValue := range cacheBookings {
		cacheBooking := cachedValue.(*models.Booking)
		// Only add if it doesn't exist in our map (to avoid duplicates)
		if _, exists := bookingMap[cacheBooking.ID]; !exists {
			bookingMap[cacheBooking.ID] = cacheBooking
		}
	}

	// Convert map back to slice
	mergedBookings := make([]*models.Booking, 0, len(bookingMap))
	for _, booking := range bookingMap {
		mergedBookings = append(mergedBookings, booking)
	}

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
	// Get booking
	booking, err := uc.GetBookingByID(ctx, id)
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

	// Update in cache
	cacheKey := fmt.Sprintf("booking:%d", id)
	uc.cache.Delete(cacheKey)

	return updatedBooking, nil
}

// checkCredit simulates a credit check for high-value bookings
func (uc *BookingUseCaseImpl) checkCredit(ctx context.Context, booking *models.Booking) {
	// Simulate some processing time
	time.Sleep(2 * time.Second)

	// Random credit check result (70% success rate)
	rand.Seed(time.Now().UnixNano())
	status := models.BookingStatusConfirmed
	if rand.Float64() < 0.3 { // 30% chance of rejection
		status = models.BookingStatusRejected
	}

	// Update booking status
	booking.Status = status
	booking.UpdatedAt = time.Now()

	// Update in repository
	updatedBooking, err := uc.repo.Update(ctx, booking)
	if err != nil {
		log.Printf("Error updating booking after credit check: %v", err)
		return
	}

	// Update in cache
	cacheKey := fmt.Sprintf("booking:%d", booking.ID)
	uc.cache.Set(cacheKey, updatedBooking)

	log.Printf("Credit check completed for booking %d. Status: %s", booking.ID, status)
}

// checkExpiredBookings runs as a background task to cancel expired bookings
func (uc *BookingUseCaseImpl) checkExpiredBookings() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("Running expired bookings check...")
		ctx := context.Background()

		// Get all bookings from repository
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
				updatedBooking, err := uc.repo.Update(ctx, booking)
				if err != nil {
					log.Printf("Error updating expired booking %d: %v", booking.ID, err)
					continue
				}

				// Update in cache
				cacheKey := fmt.Sprintf("booking:%d", booking.ID)
				uc.cache.Set(cacheKey, updatedBooking)

				expiredCount++
			}
		}

		if expiredCount > 0 {
			log.Printf("Auto-canceled %d expired bookings", expiredCount)
		}
	}
}
