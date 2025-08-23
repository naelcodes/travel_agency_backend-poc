package domains

import (
	"errors"
	"fmt"
	"time"

	CustomErrors "neema.co.za/rest/utils/errors"
	"neema.co.za/rest/utils/helpers"
	"neema.co.za/rest/utils/models"
)

type InvoiceDomain struct {
	invoice *models.Invoice
	errors  error
}

func NewInvoiceDomain(invoice *models.Invoice) *InvoiceDomain {
	domain := &InvoiceDomain{invoice: invoice}
	return domain
}

func (domain *InvoiceDomain) SetDefaults() {
	domain.invoice.Status = "unpaid"
	domain.invoice.Tag = "3"
	domain.invoice.BaseAmount = domain.invoice.Amount
	domain.invoice.NetAmount = domain.invoice.Amount
	domain.invoice.CreditApply = 0
	domain.invoice.Balance = domain.invoice.Amount
}

func (domain *InvoiceDomain) GetInvoice() *models.Invoice {
	return domain.invoice
}

func (domain *InvoiceDomain) CheckDates() error {
	creationDate, _ := time.Parse("2006-01-02", domain.invoice.CreationDate)
	dueDate, _ := time.Parse("2006-01-02", domain.invoice.DueDate)
	if dueDate.Before(creationDate) {
		return CustomErrors.DomainError(fmt.Errorf("due date can't be less than creation date"))
	}
	return nil
}
func (domain *InvoiceDomain) ApplyImputation(imputedAmount float64) error {

	domain.invoice.CreditApply += helpers.RoundDecimalPlaces(imputedAmount, 2)
	err := domain.CalculateBalance()

	if err != nil {
		return err
	}
	return nil
}

func (domain *InvoiceDomain) CalculateBalance() error {

	if domain.invoice.CreditApply > domain.invoice.Amount {
		return CustomErrors.DomainError(fmt.Errorf("the balance of an invoice can't be less than 0. credit_apply can't be greater than  invoice amount"))
	}
	domain.invoice.Balance = helpers.RoundDecimalPlaces(domain.invoice.Amount-domain.invoice.CreditApply, 2)
	domain.UpdateStatus()

	return nil
}

func (domain *InvoiceDomain) UpdateStatus() {

	if domain.invoice.CreditApply == domain.invoice.Amount && domain.invoice.Balance == 0 {
		domain.invoice.Status = "paid"
	} else {
		domain.invoice.Status = "unpaid"
	}
}

func (domain *InvoiceDomain) Validate() error {
	if domain.invoice.CreditApply > domain.invoice.Amount {
		domain.errors = errors.Join(domain.errors, fmt.Errorf("invoice.credit_apply can't be greater than  invoice.amount"))
	}
	if domain.invoice.Balance < 0 {
		domain.errors = errors.Join(domain.errors, fmt.Errorf("invoice.balance can't be less than 0"))
	}

	if domain.invoice.Balance != helpers.RoundDecimalPlaces(domain.invoice.Amount-domain.invoice.CreditApply, 2) {
		domain.errors = errors.Join(domain.errors, fmt.Errorf("invoice.balance is not equal to the difference between invoice.amount and invoice.credit_apply"))
	}

	if domain.invoice.CreditApply < 0 {
		domain.errors = errors.Join(domain.errors, fmt.Errorf("invoice.credit_apply can't be less than 0"))
	}

	err := domain.CheckDates()
	if err != nil {
		domain.errors = errors.Join(domain.errors, err)
	}

	if domain.errors != nil {
		return CustomErrors.DomainError(domain.errors)
	}

	return nil
}
