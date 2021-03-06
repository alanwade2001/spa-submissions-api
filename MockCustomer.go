package main

// MockCustomerService s
type MockCustomerService struct {
}

// NewMockCustomerService f
func NewMockCustomerService() CustomerAPI {
	return MockCustomerService{}
}

// Find f
func (cs MockCustomerService) Find(user User) (*Customer, error) {

	return &Customer{ID: "1", Name: "Corporation ABC", InitiatingPartyID: "112233"}, nil
}
