package payloads

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"neema.co.za/rest/utils/models"
)

type CreateInvoicePayload struct {
	models.Invoice
	TravelItemIds []int `json:"travelItemIds"`
}

func (c CreateInvoicePayload) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.IdCustomer, validation.Required),
		validation.Field(&c.CreationDate, validation.Required, validation.Date("2006-01-02")),
		validation.Field(&c.DueDate, validation.Required, validation.Date("2006-01-02")),
		validation.Field(&c.TravelItemIds, validation.Required),
	)
}
