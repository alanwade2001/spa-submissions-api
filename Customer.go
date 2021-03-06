package main

// CustomerService s
type CustomerService struct {
}

// NewCustomerService f
func NewCustomerService() CustomerAPI {
	return CustomerService{}
}

// Find f
func (cs CustomerService) Find(user User) (*Customer, error) {

	return &Customer{ID: "1", Name: "Corporation ABC", InitiatingPartyID: "112233"}, nil
}
