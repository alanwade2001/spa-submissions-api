package main

import (
	"fmt"

	types "github.com/alanwade2001/spa-common"
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
func (cs CustomerService) Find(user types.UserReference) (customer *types.CustomerReference, err error) {

	// Create a resty client
	client := resty.New()
	customerURITemplate := viper.GetString("CUSTOMER_URI_TEMPLATE")
	customerURI := fmt.Sprintf(customerURITemplate, user.Email)
	var resp *resty.Response

	if resp, err = client.R().SetResult(&types.CustomerReference{}).Get(customerURI); err != nil {
		return nil, err
	}

	customer = resp.Result().(*types.CustomerReference)

	return customer, nil
}
