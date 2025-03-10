# Fiber Booking System

A booking management system built with Go and Fiber framework, using Clean Architecture principles.

## Features

- Create, view, and cancel bookings
- Cache-first data retrieval
- Asynchronous credit checking for high-value bookings (>50,000)
- Background task for auto-canceling expired bookings
- Swagger API documentation
- In-memory caching
- Comprehensive error handling
- Sorting and filtering options

## Project Structure

```
fiber-booking-system/
|— cmd/
|   |— main.go         # Entry point
|— router/
|   |— router.go       # API route definitions
|— middleware/
|   |— auth.go         # Authentication middleware
|   |— logging.go      # Logging middleware
|— dto/
|   |— booking.go      # Data Transfer Objects
|— handler/
|   |— booking_handler.go # Controller layer
|— usecase/
|   |— booking_usecase.go # Business logic layer
|— repository/
|   |— booking_repo.go    # Repository layer (Mock)
|— models/
|   |— booking.go      # Data models
|— utils/
|   |— cache.go        # Cache implementation
|   |— hash.go         # Hash and sorting utilities
|— docs/
|   |— swagger.yaml    # Swagger API documentation
|— go.mod
|— go.sum
|— README.md
```

## Requirements

- Go 1.16 or higher
- Fiber v2

## Quick Start

1. Clone the repository:
   ```
   git clone https://github.com/your-username/fiber-booking-system.git
   cd fiber-booking-system
   ```

2. Install dependencies:
   ```
   go mod download
   ```

3. Run the application:
   ```
   go run cmd/main.go
   ```

4. Access the API at `http://localhost:3000/api`
5. View the Swagger documentation at `http://localhost:3000/swagger/index.html`

## API Endpoints

- `POST /api/bookings` - Create a new booking
- `GET /api/bookings/{id}` - Get a booking by ID
- `GET /api/bookings` - Get all bookings
  - Query Parameters:
    - `sort` - Sort bookings by 'price' or 'date'
    - `high-value` - Filter high-value bookings (price > 50,000)
- `DELETE /api/bookings/{id}` - Cancel a booking

## Implementation Details

### Cache System
- In-memory cache implementation
- Bookings are stored in cache for quick retrieval
- Cache is updated when bookings are created, modified, or deleted

### Background Tasks
- High-value bookings (>50,000) trigger asynchronous credit checks
- A background task runs every minute to auto-cancel bookings that have been in 'pending' status for more than 5 minutes

### Mock Repository
- The repository layer uses a mock implementation
- Default bookings with IDs 1-10 are pre-populated
- Changes are stored in memory during the application's lifetime

## Testing

Run the tests:
```
go test ./...
```

## License

MIT