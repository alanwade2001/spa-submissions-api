package main

import (
	types "github.com/alanwade2001/spa-common"
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
func (us UserService) Find(ctx *gin.Context) (*types.UserReference, error) {

	email := ctx.GetHeader("Email")

	if email == "" {
		return nil, nil
	}

	return &types.UserReference{Email: email}, nil
}
