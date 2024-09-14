// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"
	domain "spacet/internal/domain"

	mock "github.com/stretchr/testify/mock"
)

// SpaceXAPIQueries is an autogenerated mock type for the SpaceXAPIQueries type
type SpaceXAPIQueries struct {
	mock.Mock
}

type SpaceXAPIQueries_Expecter struct {
	mock *mock.Mock
}

func (_m *SpaceXAPIQueries) EXPECT() *SpaceXAPIQueries_Expecter {
	return &SpaceXAPIQueries_Expecter{mock: &_m.Mock}
}

// GetLaunchPads provides a mock function with given fields: ctx
func (_m *SpaceXAPIQueries) GetLaunchPads(ctx context.Context) ([]*domain.LaunchPad, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetLaunchPads")
	}

	var r0 []*domain.LaunchPad
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]*domain.LaunchPad, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []*domain.LaunchPad); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.LaunchPad)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SpaceXAPIQueries_GetLaunchPads_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetLaunchPads'
type SpaceXAPIQueries_GetLaunchPads_Call struct {
	*mock.Call
}

// GetLaunchPads is a helper method to define mock.On call
//   - ctx context.Context
func (_e *SpaceXAPIQueries_Expecter) GetLaunchPads(ctx interface{}) *SpaceXAPIQueries_GetLaunchPads_Call {
	return &SpaceXAPIQueries_GetLaunchPads_Call{Call: _e.mock.On("GetLaunchPads", ctx)}
}

func (_c *SpaceXAPIQueries_GetLaunchPads_Call) Run(run func(ctx context.Context)) *SpaceXAPIQueries_GetLaunchPads_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *SpaceXAPIQueries_GetLaunchPads_Call) Return(ret []*domain.LaunchPad, _a1 error) *SpaceXAPIQueries_GetLaunchPads_Call {
	_c.Call.Return(ret, _a1)
	return _c
}

func (_c *SpaceXAPIQueries_GetLaunchPads_Call) RunAndReturn(run func(context.Context) ([]*domain.LaunchPad, error)) *SpaceXAPIQueries_GetLaunchPads_Call {
	_c.Call.Return(run)
	return _c
}

// GetUpcomingLaunches provides a mock function with given fields: ctx
func (_m *SpaceXAPIQueries) GetUpcomingLaunches(ctx context.Context) ([]*domain.Launch, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetUpcomingLaunches")
	}

	var r0 []*domain.Launch
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]*domain.Launch, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []*domain.Launch); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.Launch)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SpaceXAPIQueries_GetUpcomingLaunches_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUpcomingLaunches'
type SpaceXAPIQueries_GetUpcomingLaunches_Call struct {
	*mock.Call
}

// GetUpcomingLaunches is a helper method to define mock.On call
//   - ctx context.Context
func (_e *SpaceXAPIQueries_Expecter) GetUpcomingLaunches(ctx interface{}) *SpaceXAPIQueries_GetUpcomingLaunches_Call {
	return &SpaceXAPIQueries_GetUpcomingLaunches_Call{Call: _e.mock.On("GetUpcomingLaunches", ctx)}
}

func (_c *SpaceXAPIQueries_GetUpcomingLaunches_Call) Run(run func(ctx context.Context)) *SpaceXAPIQueries_GetUpcomingLaunches_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *SpaceXAPIQueries_GetUpcomingLaunches_Call) Return(_a0 []*domain.Launch, _a1 error) *SpaceXAPIQueries_GetUpcomingLaunches_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SpaceXAPIQueries_GetUpcomingLaunches_Call) RunAndReturn(run func(context.Context) ([]*domain.Launch, error)) *SpaceXAPIQueries_GetUpcomingLaunches_Call {
	_c.Call.Return(run)
	return _c
}

// NewSpaceXAPIQueries creates a new instance of SpaceXAPIQueries. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSpaceXAPIQueries(t interface {
	mock.TestingT
	Cleanup(func())
}) *SpaceXAPIQueries {
	mock := &SpaceXAPIQueries{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
