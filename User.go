package main

import (
	"github.com/gin-gonic/gin"
)

// UserService s
type UserService struct {
}

// NewUserService f
func NewUserService() UserAPI {
	return UserService{}
}

// Find f
func (us UserService) Find(ctx *gin.Context) (*User, error) {

	email := ctx.GetHeader("Email")

	if email == "" {
		return nil, nil
	}

	return &User{Email: email}, nil
}
