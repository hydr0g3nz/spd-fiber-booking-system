# Fiber Booking System

A robust booking management system built with Go and Fiber framework, following Clean Architecture principles.

## Features

- **RESTful API Endpoints**: Create, view, and cancel bookings through a clean API interface
- **Clean Architecture**: Separation of concerns with layered design (handlers, use cases, repositories)
- **Cache-First Strategy**: Optimized performance with in-memory caching
- **Asynchronous Processing**: Background processing for high-value bookings (>50,000)
- **Auto-Cancellation**: Background task for auto-canceling expired bookings
- **Type-Safe Enums**: Strongly-typed booking status using custom Go enums
- **Comprehensive Testing**: Unit tests for all layers with high coverage
- **API Documentation**: Swagger OpenAPI documentation
- **Filter & Sort Options**: Advanced querying capabilities
- **Error Handling**: Robust error handling throughout the application
- **Authentication Middleware**: API key validation for securing endpoints

## Project Architecture

The project follows Robert C. Martin's Clean Architecture principles with clear separation of concerns:

```
fiber-booking-system/
|— cmd/                 # Application entry point
|— handler/             # HTTP request handlers (controller layer)
|— usecase/             # Business logic layer
|— repository/          # Data access layer
|— models/              # Domain models and entities
|— dto/                 # Data Transfer Objects
|— middleware/          # HTTP middleware
|— router/              # API route definitions
|— utils/               # Helper utilities
|— docs/                # Swagger documentation
|— mocks/               # Mock implementations for testing
```

### Clean Architecture Flow

1. External requests enter through the **handler layer**
2. Handlers validate input and call the appropriate **use cases**
3. Use cases implement business logic and interact with **repositories**
4. Repositories handle data access and storage
5. Data flows through the layers using **models** and **DTOs**

## Special Features

### Authentication Middleware

The system implements an API key-based authentication middleware:

```go
// Auth middleware for authentication
func Auth() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Check for API key in headers
        apiKey := c.Get("X-API-Key")
        
        // For demo, we'll accept any API key that's at least 10 characters
        if len(apiKey) < 10 {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Invalid API Key",
            })
        }

        // Authentication successful
        return c.Next()
    }
}
```

All API endpoints require a valid API key (at least 10 characters long) to be provided in the `X-API-Key` header.

### Type-Safe Enum Implementation

The project uses a type-safe enum pattern for booking status:

```go
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
```

This approach provides type safety, IDE autocompletion, and prevents invalid status values.

### Asynchronous Credit Checking

High-value bookings (above 50,000) trigger asynchronous credit checks:

```go
// For high-value bookings, run credit check in background
if newBooking.Price > 50000 {
    go uc.checkCredit(ctx, newBooking)
}
```

### Cache-First Data Access Strategy

The application uses a cache-first strategy for optimal performance:

```go
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
```

### Background Task for Auto-Cancellation

A background task runs periodically to auto-cancel expired bookings:

```go
// If booking is pending for more than 5 minutes, mark as canceled
if booking.Status == models.BookingStatusPending && now.Sub(booking.CreatedAt) > 5*time.Minute {
    booking.Status = models.BookingStatusCanceled
    booking.UpdatedAt = now
    
    // Update in repository
    updatedBooking, err := uc.repo.Update(ctx, booking)
    // ...
}
```

## Requirements

- Go 1.16 or higher
- Fiber v2.0 or higher
- Swagger tools for API documentation

## Getting Started

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/fiber-booking-system.git
   cd fiber-booking-system
   ```

2. Install dependencies:
   ```bash
   make deps
   # or
   go mod download
   ```

### Running the Application

Use the Makefile to run the application:

```bash
make run
```

Or run directly with Go:

```bash
go run cmd/main.go
```

The server will start on `http://localhost:3000`.

### Building the Application

```bash
make build
```

This will create a binary in the `./build` directory.

### Running Tests

Run all tests:

```bash
make test
# or
go test ./...
```

Run tests with coverage:

```bash
make test-coverage
```

### Generate Swagger Documentation

```bash
make docs
# or
swag init -g cmd/main.go -o docs
```

## API Documentation

Swagger API documentation is available at:
```
http://localhost:3000/swagger/index.html
```

### API Endpoints

- `POST /api/bookings` - Create a new booking
- `GET /api/bookings/{id}` - Get a booking by ID
- `GET /api/bookings` - Get all bookings
  - Query Parameters:
    - `sort` - Sort bookings by 'price' or 'date'
    - `high-value` - Filter high-value bookings (price > 50,000)
- `DELETE /api/bookings/{id}` - Cancel a booking

### Authentication

All API endpoints require authentication using an API key:

- Header: `X-API-Key`
- Format: Any string of at least 10 characters

Example:
```
X-API-Key: abcdef1234567890
```

## Implementation Details

### Cache System
- In-memory cache implementation with thread-safe operations
- Bookings are stored in cache for quick retrieval
- Cache is updated when bookings are created, modified, or deleted

### Background Tasks
- High-value bookings (>50,000) trigger asynchronous credit checks
- A background task runs every minute to auto-cancel bookings that have been in 'pending' status for more than 5 minutes

### Mock Repository
- The repository layer uses a mock implementation for demonstration
- Default bookings with IDs 1-10 are pre-populated
- Changes are stored in memory during the application's lifetime

## Development Workflow

1. Make changes to the code
2. Run tests: `make test`
3. Regenerate Swagger docs if needed: `make docs`
4. Run the application: `make run`
5. Check the API at http://localhost:3000/api
6. View Swagger documentation at http://localhost:3000/swagger/index.html

## License

MIT