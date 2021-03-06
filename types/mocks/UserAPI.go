// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	gin "github.com/gin-gonic/gin"
	mock "github.com/stretchr/testify/mock"

	submission "github.com/alanwade2001/spa-submissions-api/models/generated/submission"
)

// UserAPI is an autogenerated mock type for the UserAPI type
type UserAPI struct {
	mock.Mock
}

// Find provides a mock function with given fields: _a0
func (_m *UserAPI) Find(_a0 *gin.Context) (*submission.UserReference, error) {
	ret := _m.Called(_a0)

	var r0 *submission.UserReference
	if rf, ok := ret.Get(0).(func(*gin.Context) *submission.UserReference); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*submission.UserReference)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*gin.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
