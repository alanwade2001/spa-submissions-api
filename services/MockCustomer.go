package services

import (
	"github.com/alanwade2001/spa-submissions-api/models/generated/submission"
	"github.com/alanwade2001/spa-submissions-api/types"
)

// MockCustomerService s
type MockCustomerService struct {
}

// NewMockCustomerService f
func NewMockCustomerService() types.CustomerAPI {
	return MockCustomerService{}
}

// Find f
func (cs MockCustomerService) Find(user submission.UserReference) (*submission.CustomerReference, error) {

	return &submission.CustomerReference{
		Id:   "1",
		Name: "Corporation ABC",
	}, nil
}
