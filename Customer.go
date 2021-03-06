package main

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

// CustomerService s
type CustomerService struct {
}

// NewCustomerService f
func NewCustomerService() CustomerAPI {
	return CustomerService{}
}

// Find f
func (cs CustomerService) Find(user User) (customer *Customer, err error) {

	// Create a resty client
	client := resty.New()
	customerURITemplate := viper.GetString("CUSTOMER_URI_TEMPLATE")
	customerURI := fmt.Sprintf(customerURITemplate, user.Email)
	var resp *resty.Response

	if resp, err = client.R().SetResult(&Customer{}).Get(customerURI); err != nil {
		return nil, err
	}

	customer = resp.Result().(*Customer)

	return customer, nil
}
