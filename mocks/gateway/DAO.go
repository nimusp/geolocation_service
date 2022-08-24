// Code generated by mockery v2.14.0. DO NOT EDIT.

package gateway

import (
	context "context"

	gateway "github.com/nimusp/geolocation_service/internal/gateway"
	mock "github.com/stretchr/testify/mock"
)

// DAO is an autogenerated mock type for the DAO type
type DAO struct {
	mock.Mock
}

type DAO_Expecter struct {
	mock *mock.Mock
}

func (_m *DAO) EXPECT() *DAO_Expecter {
	return &DAO_Expecter{mock: &_m.Mock}
}

// GetByIP provides a mock function with given fields: _a0, _a1
func (_m *DAO) GetByIP(_a0 context.Context, _a1 string) (*gateway.GeoLocation, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *gateway.GeoLocation
	if rf, ok := ret.Get(0).(func(context.Context, string) *gateway.GeoLocation); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gateway.GeoLocation)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DAO_GetByIP_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByIP'
type DAO_GetByIP_Call struct {
	*mock.Call
}

// GetByIP is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 string
func (_e *DAO_Expecter) GetByIP(_a0 interface{}, _a1 interface{}) *DAO_GetByIP_Call {
	return &DAO_GetByIP_Call{Call: _e.mock.On("GetByIP", _a0, _a1)}
}

func (_c *DAO_GetByIP_Call) Run(run func(_a0 context.Context, _a1 string)) *DAO_GetByIP_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *DAO_GetByIP_Call) Return(_a0 *gateway.GeoLocation, _a1 error) *DAO_GetByIP_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

type mockConstructorTestingTNewDAO interface {
	mock.TestingT
	Cleanup(func())
}

// NewDAO creates a new instance of DAO. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewDAO(t mockConstructorTestingTNewDAO) *DAO {
	mock := &DAO{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
