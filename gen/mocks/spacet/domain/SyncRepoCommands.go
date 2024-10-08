// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// SyncRepoCommands is an autogenerated mock type for the SyncRepoCommands type
type SyncRepoCommands struct {
	mock.Mock
}

type SyncRepoCommands_Expecter struct {
	mock *mock.Mock
}

func (_m *SyncRepoCommands) EXPECT() *SyncRepoCommands_Expecter {
	return &SyncRepoCommands_Expecter{mock: &_m.Mock}
}

// GetLastSyncTimestamp provides a mock function with given fields: ctx, resourceName
func (_m *SyncRepoCommands) GetLastSyncTimestamp(ctx context.Context, resourceName string) (time.Time, error) {
	ret := _m.Called(ctx, resourceName)

	if len(ret) == 0 {
		panic("no return value specified for GetLastSyncTimestamp")
	}

	var r0 time.Time
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (time.Time, error)); ok {
		return rf(ctx, resourceName)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) time.Time); ok {
		r0 = rf(ctx, resourceName)
	} else {
		r0 = ret.Get(0).(time.Time)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, resourceName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SyncRepoCommands_GetLastSyncTimestamp_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetLastSyncTimestamp'
type SyncRepoCommands_GetLastSyncTimestamp_Call struct {
	*mock.Call
}

// GetLastSyncTimestamp is a helper method to define mock.On call
//   - ctx context.Context
//   - resourceName string
func (_e *SyncRepoCommands_Expecter) GetLastSyncTimestamp(ctx interface{}, resourceName interface{}) *SyncRepoCommands_GetLastSyncTimestamp_Call {
	return &SyncRepoCommands_GetLastSyncTimestamp_Call{Call: _e.mock.On("GetLastSyncTimestamp", ctx, resourceName)}
}

func (_c *SyncRepoCommands_GetLastSyncTimestamp_Call) Run(run func(ctx context.Context, resourceName string)) *SyncRepoCommands_GetLastSyncTimestamp_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *SyncRepoCommands_GetLastSyncTimestamp_Call) Return(_a0 time.Time, _a1 error) *SyncRepoCommands_GetLastSyncTimestamp_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SyncRepoCommands_GetLastSyncTimestamp_Call) RunAndReturn(run func(context.Context, string) (time.Time, error)) *SyncRepoCommands_GetLastSyncTimestamp_Call {
	_c.Call.Return(run)
	return _c
}

// ReleaseDistributedLock provides a mock function with given fields: ctx, key
func (_m *SyncRepoCommands) ReleaseDistributedLock(ctx context.Context, key uint32) error {
	ret := _m.Called(ctx, key)

	if len(ret) == 0 {
		panic("no return value specified for ReleaseDistributedLock")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32) error); ok {
		r0 = rf(ctx, key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SyncRepoCommands_ReleaseDistributedLock_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ReleaseDistributedLock'
type SyncRepoCommands_ReleaseDistributedLock_Call struct {
	*mock.Call
}

// ReleaseDistributedLock is a helper method to define mock.On call
//   - ctx context.Context
//   - key uint32
func (_e *SyncRepoCommands_Expecter) ReleaseDistributedLock(ctx interface{}, key interface{}) *SyncRepoCommands_ReleaseDistributedLock_Call {
	return &SyncRepoCommands_ReleaseDistributedLock_Call{Call: _e.mock.On("ReleaseDistributedLock", ctx, key)}
}

func (_c *SyncRepoCommands_ReleaseDistributedLock_Call) Run(run func(ctx context.Context, key uint32)) *SyncRepoCommands_ReleaseDistributedLock_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint32))
	})
	return _c
}

func (_c *SyncRepoCommands_ReleaseDistributedLock_Call) Return(_a0 error) *SyncRepoCommands_ReleaseDistributedLock_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *SyncRepoCommands_ReleaseDistributedLock_Call) RunAndReturn(run func(context.Context, uint32) error) *SyncRepoCommands_ReleaseDistributedLock_Call {
	_c.Call.Return(run)
	return _c
}

// TryDistributedLock provides a mock function with given fields: ctx, key
func (_m *SyncRepoCommands) TryDistributedLock(ctx context.Context, key uint32) (bool, error) {
	ret := _m.Called(ctx, key)

	if len(ret) == 0 {
		panic("no return value specified for TryDistributedLock")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32) (bool, error)); ok {
		return rf(ctx, key)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32) bool); ok {
		r0 = rf(ctx, key)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32) error); ok {
		r1 = rf(ctx, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SyncRepoCommands_TryDistributedLock_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'TryDistributedLock'
type SyncRepoCommands_TryDistributedLock_Call struct {
	*mock.Call
}

// TryDistributedLock is a helper method to define mock.On call
//   - ctx context.Context
//   - key uint32
func (_e *SyncRepoCommands_Expecter) TryDistributedLock(ctx interface{}, key interface{}) *SyncRepoCommands_TryDistributedLock_Call {
	return &SyncRepoCommands_TryDistributedLock_Call{Call: _e.mock.On("TryDistributedLock", ctx, key)}
}

func (_c *SyncRepoCommands_TryDistributedLock_Call) Run(run func(ctx context.Context, key uint32)) *SyncRepoCommands_TryDistributedLock_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint32))
	})
	return _c
}

func (_c *SyncRepoCommands_TryDistributedLock_Call) Return(_a0 bool, _a1 error) *SyncRepoCommands_TryDistributedLock_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SyncRepoCommands_TryDistributedLock_Call) RunAndReturn(run func(context.Context, uint32) (bool, error)) *SyncRepoCommands_TryDistributedLock_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateLastSyncTimestamp provides a mock function with given fields: ctx, resourceName, newTimestamp
func (_m *SyncRepoCommands) UpdateLastSyncTimestamp(ctx context.Context, resourceName string, newTimestamp time.Time) error {
	ret := _m.Called(ctx, resourceName, newTimestamp)

	if len(ret) == 0 {
		panic("no return value specified for UpdateLastSyncTimestamp")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, time.Time) error); ok {
		r0 = rf(ctx, resourceName, newTimestamp)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SyncRepoCommands_UpdateLastSyncTimestamp_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateLastSyncTimestamp'
type SyncRepoCommands_UpdateLastSyncTimestamp_Call struct {
	*mock.Call
}

// UpdateLastSyncTimestamp is a helper method to define mock.On call
//   - ctx context.Context
//   - resourceName string
//   - newTimestamp time.Time
func (_e *SyncRepoCommands_Expecter) UpdateLastSyncTimestamp(ctx interface{}, resourceName interface{}, newTimestamp interface{}) *SyncRepoCommands_UpdateLastSyncTimestamp_Call {
	return &SyncRepoCommands_UpdateLastSyncTimestamp_Call{Call: _e.mock.On("UpdateLastSyncTimestamp", ctx, resourceName, newTimestamp)}
}

func (_c *SyncRepoCommands_UpdateLastSyncTimestamp_Call) Run(run func(ctx context.Context, resourceName string, newTimestamp time.Time)) *SyncRepoCommands_UpdateLastSyncTimestamp_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(time.Time))
	})
	return _c
}

func (_c *SyncRepoCommands_UpdateLastSyncTimestamp_Call) Return(_a0 error) *SyncRepoCommands_UpdateLastSyncTimestamp_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *SyncRepoCommands_UpdateLastSyncTimestamp_Call) RunAndReturn(run func(context.Context, string, time.Time) error) *SyncRepoCommands_UpdateLastSyncTimestamp_Call {
	_c.Call.Return(run)
	return _c
}

// NewSyncRepoCommands creates a new instance of SyncRepoCommands. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSyncRepoCommands(t interface {
	mock.TestingT
	Cleanup(func())
}) *SyncRepoCommands {
	mock := &SyncRepoCommands{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
