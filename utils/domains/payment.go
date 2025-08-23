package domains

import (
	"errors"
	"fmt"

	CustomErrors "neema.co.za/rest/utils/errors"
	"neema.co.za/rest/utils/helpers"
	"neema.co.za/rest/utils/models"
)

type PaymentDomain struct {
	payment *models.Payment
	errors  error
}

func NewPaymentDomain(payment *models.Payment) *PaymentDomain {
	domain := &PaymentDomain{payment: payment}
	return domain

}

func (domain *PaymentDomain) SetDefaults() {
	domain.payment.PaymentDate = helpers.GetCurrentDate()
	domain.payment.BaseAmount = domain.payment.Amount
	domain.payment.Status = "open"
	domain.payment.UsedAmount = 0
	domain.payment.Type = "customer_payment"
	domain.payment.IdCurrency = 550
	domain.payment.IdChartOfAccounts = 39
	domain.payment.Balance = domain.payment.Amount
	domain.payment.Tag = "3"
}

func (domain *PaymentDomain) GetPayment() *models.Payment {
	return domain.payment
}

func (domain *PaymentDomain) calculateBalance() error {

	if domain.payment.UsedAmount > domain.payment.Amount {
		return CustomErrors.DomainError(errors.New("payment balance can't be less than 0"))
	}

	domain.payment.Balance = helpers.RoundDecimalPlaces(domain.payment.Amount-domain.payment.UsedAmount, 2)
	domain.updateStatus()

	return nil
}

func (domain *PaymentDomain) updateStatus() {

	if domain.payment.UsedAmount == domain.payment.Amount && domain.payment.Balance == 0 {
		domain.payment.Status = "used"
	} else {
		domain.payment.Status = "open"
	}
}

func (domain *PaymentDomain) AllocateAmount(oldImputationAmount float64, newImputationAmount float64) error {

	domain.payment.UsedAmount = domain.payment.UsedAmount - helpers.RoundDecimalPlaces(oldImputationAmount, 2)

	if domain.payment.UsedAmount+newImputationAmount > domain.payment.Amount {
		return CustomErrors.DomainError(fmt.Errorf("allocated(used) amount on payment %v can't be greater than the payment amount", domain.payment.PaymentNumber))
	}

	domain.payment.UsedAmount = domain.payment.UsedAmount + helpers.RoundDecimalPlaces(newImputationAmount, 2)
	err := domain.calculateBalance()

	if err != nil {
		return err
	}

	return nil
}

func (domain *PaymentDomain) Validate() error {

	if domain.payment.Amount < 0 {
		domain.errors = errors.Join(domain.errors, fmt.Errorf("payment.amount can't be less than 0"))
	}

	if domain.payment.Balance < 0 {
		domain.errors = errors.Join(domain.errors, fmt.Errorf("payment.balance can't be less than 0"))
	}

	if domain.payment.UsedAmount < 0 {
		domain.errors = errors.Join(domain.errors, fmt.Errorf("payment.usedAmount can't be less than 0"))
	}

	if domain.payment.Balance != helpers.RoundDecimalPlaces(domain.payment.Amount-domain.payment.UsedAmount, 2) {
		domain.errors = errors.Join(domain.errors, fmt.Errorf("payment.balance (%v) must be equal to the difference between payment.amount (%v) and payment.usedAmount (%v)", domain.payment.Balance, float64(domain.payment.Amount), float64(domain.payment.UsedAmount)))
	}

	if domain.payment.UsedAmount > domain.payment.Amount {
		domain.errors = errors.Join(domain.errors, fmt.Errorf("payment.usedAmount can't be greater than payment.amount"))
	}

	if domain.errors != nil {
		return CustomErrors.DomainError(domain.errors)
	}
	return nil

}
