// Code generated by mockery v2.53.1. DO NOT EDIT.

package mocks

import (
	context "context"

	dto "github.com/hydr0g3nz/spd-fiber-booking-system/dto"
	mock "github.com/stretchr/testify/mock"

	models "github.com/hydr0g3nz/spd-fiber-booking-system/models"
)

// BookingUseCase is an autogenerated mock type for the BookingUseCase type
type BookingUseCase struct {
	mock.Mock
}

// CancelBooking provides a mock function with given fields: ctx, id
func (_m *BookingUseCase) CancelBooking(ctx context.Context, id int64) (*models.Booking, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for CancelBooking")
	}

	var r0 *models.Booking
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (*models.Booking, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) *models.Booking); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Booking)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateBooking provides a mock function with given fields: ctx, req
func (_m *BookingUseCase) CreateBooking(ctx context.Context, req *dto.CreateBookingRequest) (*models.Booking, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for CreateBooking")
	}

	var r0 *models.Booking
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *dto.CreateBookingRequest) (*models.Booking, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *dto.CreateBookingRequest) *models.Booking); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Booking)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *dto.CreateBookingRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllBookings provides a mock function with given fields: ctx, params
func (_m *BookingUseCase) GetAllBookings(ctx context.Context, params *dto.BookingsQueryParams) ([]*models.Booking, error) {
	ret := _m.Called(ctx, params)

	if len(ret) == 0 {
		panic("no return value specified for GetAllBookings")
	}

	var r0 []*models.Booking
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *dto.BookingsQueryParams) ([]*models.Booking, error)); ok {
		return rf(ctx, params)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *dto.BookingsQueryParams) []*models.Booking); ok {
		r0 = rf(ctx, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Booking)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *dto.BookingsQueryParams) error); ok {
		r1 = rf(ctx, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBookingByID provides a mock function with given fields: ctx, id
func (_m *BookingUseCase) GetBookingByID(ctx context.Context, id int64) (*models.Booking, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetBookingByID")
	}

	var r0 *models.Booking
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (*models.Booking, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) *models.Booking); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Booking)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewBookingUseCase creates a new instance of BookingUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewBookingUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *BookingUseCase {
	mock := &BookingUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
