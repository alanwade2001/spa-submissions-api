package services

import (
	"fmt"

	"github.com/alanwade2001/spa-submissions-api/models/generated/customer"
	"github.com/alanwade2001/spa-submissions-api/models/generated/submission"
	"github.com/alanwade2001/spa-submissions-api/types"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

// CustomerService s
type CustomerService struct {
}

// NewCustomerService f
func NewCustomerService() types.CustomerAPI {
	return CustomerService{}
}

// Find f
func (cs CustomerService) Find(user submission.UserReference) (customerRef *submission.CustomerReference, err error) {

	// Create a resty client
	client := resty.New()
	client.SetDebug(true)
	customerURITemplate := viper.GetString("CUSTOMER_URI_TEMPLATE")
	customerURI := fmt.Sprintf(customerURITemplate, user.Email)
	var resp *resty.Response

	if resp, err = client.R().SetResult(&customer.CustomerModel{}).Get(customerURI); err != nil {
		return nil, err
	}

	customerModel := resp.Result().(*customer.CustomerModel)

	customerRef = &submission.CustomerReference{
		Id:   customerModel.Id,
		Name: customerModel.Name,
	}

	return customerRef, nil
}
