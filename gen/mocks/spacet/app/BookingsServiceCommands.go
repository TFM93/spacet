// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"
	bookings "spacet/internal/app/bookings"

	domain "spacet/internal/domain"

	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// BookingsServiceCommands is an autogenerated mock type for the BookingsServiceCommands type
type BookingsServiceCommands struct {
	mock.Mock
}

type BookingsServiceCommands_Expecter struct {
	mock *mock.Mock
}

func (_m *BookingsServiceCommands) EXPECT() *BookingsServiceCommands_Expecter {
	return &BookingsServiceCommands_Expecter{mock: &_m.Mock}
}

// BookALaunch provides a mock function with given fields: ctx, req
func (_m *BookingsServiceCommands) BookALaunch(ctx context.Context, req bookings.BookALaunchReq) (domain.Ticket, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for BookALaunch")
	}

	var r0 domain.Ticket
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, bookings.BookALaunchReq) (domain.Ticket, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, bookings.BookALaunchReq) domain.Ticket); ok {
		r0 = rf(ctx, req)
	} else {
		r0 = ret.Get(0).(domain.Ticket)
	}

	if rf, ok := ret.Get(1).(func(context.Context, bookings.BookALaunchReq) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BookingsServiceCommands_BookALaunch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'BookALaunch'
type BookingsServiceCommands_BookALaunch_Call struct {
	*mock.Call
}

// BookALaunch is a helper method to define mock.On call
//   - ctx context.Context
//   - req bookings.BookALaunchReq
func (_e *BookingsServiceCommands_Expecter) BookALaunch(ctx interface{}, req interface{}) *BookingsServiceCommands_BookALaunch_Call {
	return &BookingsServiceCommands_BookALaunch_Call{Call: _e.mock.On("BookALaunch", ctx, req)}
}

func (_c *BookingsServiceCommands_BookALaunch_Call) Run(run func(ctx context.Context, req bookings.BookALaunchReq)) *BookingsServiceCommands_BookALaunch_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(bookings.BookALaunchReq))
	})
	return _c
}

func (_c *BookingsServiceCommands_BookALaunch_Call) Return(createdBooking domain.Ticket, _a1 error) *BookingsServiceCommands_BookALaunch_Call {
	_c.Call.Return(createdBooking, _a1)
	return _c
}

func (_c *BookingsServiceCommands_BookALaunch_Call) RunAndReturn(run func(context.Context, bookings.BookALaunchReq) (domain.Ticket, error)) *BookingsServiceCommands_BookALaunch_Call {
	_c.Call.Return(run)
	return _c
}

// Cancel provides a mock function with given fields: ctx, restriction
func (_m *BookingsServiceCommands) Cancel(ctx context.Context, restriction bookings.LaunchPadDateRestrictions) ([]uuid.UUID, error) {
	ret := _m.Called(ctx, restriction)

	if len(ret) == 0 {
		panic("no return value specified for Cancel")
	}

	var r0 []uuid.UUID
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, bookings.LaunchPadDateRestrictions) ([]uuid.UUID, error)); ok {
		return rf(ctx, restriction)
	}
	if rf, ok := ret.Get(0).(func(context.Context, bookings.LaunchPadDateRestrictions) []uuid.UUID); ok {
		r0 = rf(ctx, restriction)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]uuid.UUID)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, bookings.LaunchPadDateRestrictions) error); ok {
		r1 = rf(ctx, restriction)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BookingsServiceCommands_Cancel_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Cancel'
type BookingsServiceCommands_Cancel_Call struct {
	*mock.Call
}

// Cancel is a helper method to define mock.On call
//   - ctx context.Context
//   - restriction bookings.LaunchPadDateRestrictions
func (_e *BookingsServiceCommands_Expecter) Cancel(ctx interface{}, restriction interface{}) *BookingsServiceCommands_Cancel_Call {
	return &BookingsServiceCommands_Cancel_Call{Call: _e.mock.On("Cancel", ctx, restriction)}
}

func (_c *BookingsServiceCommands_Cancel_Call) Run(run func(ctx context.Context, restriction bookings.LaunchPadDateRestrictions)) *BookingsServiceCommands_Cancel_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(bookings.LaunchPadDateRestrictions))
	})
	return _c
}

func (_c *BookingsServiceCommands_Cancel_Call) Return(cancelled []uuid.UUID, err error) *BookingsServiceCommands_Cancel_Call {
	_c.Call.Return(cancelled, err)
	return _c
}

func (_c *BookingsServiceCommands_Cancel_Call) RunAndReturn(run func(context.Context, bookings.LaunchPadDateRestrictions) ([]uuid.UUID, error)) *BookingsServiceCommands_Cancel_Call {
	_c.Call.Return(run)
	return _c
}

// NewBookingsServiceCommands creates a new instance of BookingsServiceCommands. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewBookingsServiceCommands(t interface {
	mock.TestingT
	Cleanup(func())
}) *BookingsServiceCommands {
	mock := &BookingsServiceCommands{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
