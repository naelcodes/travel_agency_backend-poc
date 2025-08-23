package payloads

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"neema.co.za/rest/utils/models"
)

type CreatePaymentPayload struct {
	models.Payment
}

func (p CreatePaymentPayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.IdCustomer, validation.Required),
		validation.Field(&p.Amount, validation.Required),
		validation.Field(&p.PaymentMode, validation.Required, validation.In("cash", "check", "bank_transfer")),
	)

}
