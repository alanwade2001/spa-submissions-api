package main

import types "github.com/alanwade2001/spa-common"

// MockCustomerService s
type MockCustomerService struct {
}

// NewMockCustomerService f
func NewMockCustomerService() CustomerAPI {
	return MockCustomerService{}
}

// Find f
func (cs MockCustomerService) Find(user types.UserReference) (*types.CustomerReference, error) {

	return &types.CustomerReference{
		CustomerID: "1",
		Name:       "Corporation ABC",
		InitiatingParty: types.InitiatingPartyReference{
			InitiatingPartyID: "112233",
		},
	}, nil
}
