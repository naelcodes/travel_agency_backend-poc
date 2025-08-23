package domains

import (
	"neema.co.za/rest/utils/helpers"
	"neema.co.za/rest/utils/models"
)

type CustomerDomain struct {
	customer *models.Customer
}

func NewCustomerDomain(customer *models.Customer) *CustomerDomain {
	domain := &CustomerDomain{customer: customer}
	return domain
}

func (domain *CustomerDomain) SetDefaults() {
	domain.customer.Tag = "3"
	domain.customer.IdCountry = 40
	domain.customer.IdCurrency = 550
	domain.customer.AbKey = helpers.GenerateRandomString(15)
}

func (domain *CustomerDomain) GetCustomer() *models.Customer {
	return domain.customer
}
