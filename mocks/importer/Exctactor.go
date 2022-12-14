// Code generated by mockery v2.14.0. DO NOT EDIT.

package importer

import (
	importer "github.com/nimusp/geolocation_service/internal/importer"
	mock "github.com/stretchr/testify/mock"
)

// Exctactor is an autogenerated mock type for the Exctactor type
type Exctactor struct {
	mock.Mock
}

type Exctactor_Expecter struct {
	mock *mock.Mock
}

func (_m *Exctactor) EXPECT() *Exctactor_Expecter {
	return &Exctactor_Expecter{mock: &_m.Mock}
}

// Extract provides a mock function with given fields:
func (_m *Exctactor) Extract() (*importer.Data, error) {
	ret := _m.Called()

	var r0 *importer.Data
	if rf, ok := ret.Get(0).(func() *importer.Data); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*importer.Data)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Exctactor_Extract_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Extract'
type Exctactor_Extract_Call struct {
	*mock.Call
}

// Extract is a helper method to define mock.On call
func (_e *Exctactor_Expecter) Extract() *Exctactor_Extract_Call {
	return &Exctactor_Extract_Call{Call: _e.mock.On("Extract")}
}

func (_c *Exctactor_Extract_Call) Run(run func()) *Exctactor_Extract_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Exctactor_Extract_Call) Return(_a0 *importer.Data, _a1 error) *Exctactor_Extract_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

type mockConstructorTestingTNewExctactor interface {
	mock.TestingT
	Cleanup(func())
}

// NewExctactor creates a new instance of Exctactor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewExctactor(t mockConstructorTestingTNewExctactor) *Exctactor {
	mock := &Exctactor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
