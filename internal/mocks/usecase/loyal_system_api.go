// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	entity "github.com/vladislaoramos/gophermart/internal/entity"

	time "time"
)

// LoyalSystemAPI is an autogenerated mock type for the LoyalSystemAPI type
type LoyalSystemAPI struct {
	mock.Mock
}

type LoyalSystemAPI_Expecter struct {
	mock *mock.Mock
}

func (_m *LoyalSystemAPI) EXPECT() *LoyalSystemAPI_Expecter {
	return &LoyalSystemAPI_Expecter{mock: &_m.Mock}
}

// GetOrderInfo provides a mock function with given fields: _a0
func (_m *LoyalSystemAPI) GetOrderInfo(_a0 string) (entity.Order, time.Duration, error) {
	ret := _m.Called(_a0)

	var r0 entity.Order
	if rf, ok := ret.Get(0).(func(string) entity.Order); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(entity.Order)
	}

	var r1 time.Duration
	if rf, ok := ret.Get(1).(func(string) time.Duration); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Get(1).(time.Duration)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string) error); ok {
		r2 = rf(_a0)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// LoyalSystemAPI_GetOrderInfo_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOrderInfo'
type LoyalSystemAPI_GetOrderInfo_Call struct {
	*mock.Call
}

// GetOrderInfo is a helper method to define mock.On call
//   - _a0 string
func (_e *LoyalSystemAPI_Expecter) GetOrderInfo(_a0 interface{}) *LoyalSystemAPI_GetOrderInfo_Call {
	return &LoyalSystemAPI_GetOrderInfo_Call{Call: _e.mock.On("GetOrderInfo", _a0)}
}

func (_c *LoyalSystemAPI_GetOrderInfo_Call) Run(run func(_a0 string)) *LoyalSystemAPI_GetOrderInfo_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *LoyalSystemAPI_GetOrderInfo_Call) Return(_a0 entity.Order, _a1 time.Duration, _a2 error) *LoyalSystemAPI_GetOrderInfo_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

type mockConstructorTestingTNewLoyalSystemAPI interface {
	mock.TestingT
	Cleanup(func())
}

// NewLoyalSystemAPI creates a new instance of LoyalSystemAPI. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewLoyalSystemAPI(t mockConstructorTestingTNewLoyalSystemAPI) *LoyalSystemAPI {
	mock := &LoyalSystemAPI{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}