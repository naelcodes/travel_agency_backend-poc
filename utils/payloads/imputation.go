package payloads

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation"
)

type ImputationPayload struct {
	IdPayment     int     `json:"idPayment"`
	AmountApplied float64 `json:"amountApplied"`
}

func (i ImputationPayload) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.IdPayment, validation.Required),
		validation.Field(&i.AmountApplied, validation.By(func(value any) error {
			floatValue, ok := value.(float64)
			if !ok {
				return errors.New("validation error : imputation amount must be a numeric value")
			}

			if floatValue < 0 {
				return errors.New("validation error : imputation amount must be greater than  zero")
			}
			return nil
		})),
	)
}
