package repository

import (
	"fmt"

	CustomErrors "neema.co.za/rest/utils/errors"
	"neema.co.za/rest/utils/logger"
	"neema.co.za/rest/utils/models"
	"xorm.io/xorm"
)

func (r *Repository) GetByInvoiceId(idInvoice int) (any, error) {

	invoices := []*models.Invoice{}
	err := r.SQL("SELECT id,amount::numeric,id_customer FROM invoice WHERE id = ?", idInvoice).Find(&invoices)

	if err != nil {
		return nil, CustomErrors.RepositoryError(fmt.Errorf("error getting invoice with id(%v): %v", idInvoice, err))
	}

	if len(invoices) == 0 {
		return nil, CustomErrors.NotFoundError(fmt.Errorf("invoice with id(%v) not found", idInvoice))
	}

	imputationAmountWithRelatedPaymentsQuery := `
		SELECT
		    i.amount_apply::NUMERIC,
		    p.id AS id,
		    p.number,
		    to_char(p.date, 'yyyy-mm-dd') AS date,
		    p.balance::NUMERIC,
		    p.amount::NUMERIC AS amount
		FROM
		    invoice_payment_received AS i
		    RIGHT OUTER JOIN payment_received AS p ON i.id_payment_received = p.id
		WHERE
		    p.id_customer = ?
		    AND ((i.id_invoice = ?
		            AND i.id_payment_received IS NOT NULL)
		        OR (p.balance::NUMERIC != 0
		            AND i.id_invoice IS NULL
		            AND i.id_payment_received IS NULL))`

	imputationAmountWithRelatedPaymentRecords := make([]*struct {
		Payment       models.Payment `xorm:"extends" json:"payment"`
		AmountApplied float64        `xorm:"amount_apply" json:"amountApplied"`
	}, 0)

	err = r.SQL(imputationAmountWithRelatedPaymentsQuery, invoices[0].IdCustomer, idInvoice).Find(&imputationAmountWithRelatedPaymentRecords)

	if err != nil {
		return nil, CustomErrors.RepositoryError(fmt.Errorf("error getting imputations of invoice with id(%v) : %v", idInvoice, err))
	}

	if len(imputationAmountWithRelatedPaymentRecords) == 0 {
		return nil, CustomErrors.NotFoundError(fmt.Errorf("imputations of invoice with id(%v) not found", idInvoice))
	}

	for i := range imputationAmountWithRelatedPaymentRecords {
		imputationAmountWithRelatedPaymentRecords[i].Payment.Balance = imputationAmountWithRelatedPaymentRecords[i].Payment.Balance + imputationAmountWithRelatedPaymentRecords[i].AmountApplied
	}

	data := new(struct {
		InvoiceAmount float64 `json:"invoiceAmount"`
		Imputations   []*struct {
			Payment       models.Payment `xorm:"extends" json:"payment"`
			AmountApplied float64        `xorm:"amount_apply" json:"amountApplied"`
		} `json:"imputations"`
	})

	data.InvoiceAmount = invoices[0].Amount
	data.Imputations = imputationAmountWithRelatedPaymentRecords

	return data, nil
}

func (r *Repository) GetByPaymentIdAndInvoiceId(idPayment int, idInvoice int) (bool, *struct {
	Imputation models.Imputation `xorm:"extends"`
	Payment    models.Payment    `xorm:"extends" `
}, error) {

	results := make([]*struct {
		Imputation models.Imputation `xorm:"extends"`
		Payment    models.Payment    `xorm:"extends" `
	}, 0)
	imputationWithRelatedPaymentQuery := `
		SELECT
		    i.amount_apply::NUMERIC,
		    i.id,
		    p.amount::NUMERIC,
		    p.balance::NUMERIC,
		    p.status,
		    p.used_amount::NUMERIC,
		    p.id
		FROM
		    invoice_payment_received AS i
		    RIGHT OUTER JOIN payment_received AS p ON i.id_payment_received = p.id
		WHERE
		    i.id_invoice = ?
		    AND i.id_payment_received = ?`

	err := r.SQL(imputationWithRelatedPaymentQuery, idInvoice, idPayment).Find(&results)

	if err != nil {
		return false, nil, CustomErrors.RepositoryError(fmt.Errorf("error getting imputation with payment id(%v) and invoice id(%v): %v", idPayment, idInvoice, err))
	}

	if len(results) == 0 {
		return false, nil, nil
	}

	logger.Info(fmt.Sprintf("Found imputation with payment id(%v) and invoice id(%v) => result : %v", idPayment, idInvoice, results[0]))

	return true, results[0], nil
}

func (r *Repository) DeleteById(transaction *xorm.Session, id int) error {

	deleteCount, err := transaction.ID(id).Delete(&models.Imputation{})

	if err != nil {
		return CustomErrors.RepositoryError(fmt.Errorf("error deleting imputation with id(%v): %v", id, err))
	}

	if deleteCount == 0 {
		return CustomErrors.NotFoundError(fmt.Errorf("imputation with id(%v) not found", id))
	}
	return nil
}

func (r *Repository) Update(transaction *xorm.Session, imputation *models.Imputation) error {

	_, err := transaction.ID(imputation.Id).Update(imputation)

	if err != nil {
		return CustomErrors.RepositoryError(fmt.Errorf("error updating imputation with id(%v): %v", imputation.Id, err))
	}
	return nil
}

func (r *Repository) Save(transaction *xorm.Session, imputation *models.Imputation) error {

	_, err := transaction.Insert(imputation)

	if err != nil {
		return CustomErrors.RepositoryError(fmt.Errorf("error saving imputation: %v", err))
	}
	return nil
}
