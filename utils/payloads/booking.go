package payloads

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"neema.co.za/rest/utils/models"
)

type TravelItemPayload struct {
	models.TravelItem
}

func (t TravelItemPayload) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.Id, validation.Required),
	)

}
