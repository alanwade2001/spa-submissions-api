package services

import (
	"github.com/alanwade2001/spa-submissions-api/models/generated/submission"
	"github.com/alanwade2001/spa-submissions-api/types"
	"github.com/gin-gonic/gin"
)

// UserService s
type UserService struct {
}

// NewUserService f
func NewUserService() types.UserAPI {
	return UserService{}
}

// Find f
func (us UserService) Find(ctx *gin.Context) (*submission.UserReference, error) {

	email := ctx.GetHeader("Email")

	if email == "" {
		return nil, nil
	}

	return &submission.UserReference{Email: email}, nil
}
