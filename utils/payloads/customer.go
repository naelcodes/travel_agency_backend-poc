package payloads

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"neema.co.za/rest/utils/models"
)

type CreateCustomerPayload struct {
	models.Customer
}

func (c CreateCustomerPayload) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.CustomerName, validation.Required),
		validation.Field(&c.State, validation.Required),
		validation.Field(&c.AccountNumber, validation.Required),
		validation.Field(&c.Alias, validation.Required),
		validation.Field(&c.TmcClientNumber, validation.Required),
	)
}

type UpdateCustomerPayload struct {
	models.Customer
}

func (u UpdateCustomerPayload) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.CustomerName, validation.NilOrNotEmpty),
		validation.Field(&u.State, validation.NilOrNotEmpty),
		validation.Field(&u.AccountNumber, validation.NilOrNotEmpty),
		validation.Field(&u.Alias, validation.NilOrNotEmpty),
		validation.Field(&u.TmcClientNumber, validation.NilOrNotEmpty),
	)
}
