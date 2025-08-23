package service

import (
	"context"

	"neema.co.za/rest/utils/managers"
	"neema.co.za/rest/utils/models"
	"neema.co.za/rest/utils/types"
	"xorm.io/xorm"
)

type Imports struct {
	dependencyManager *managers.DependencyManager
}

func (i *Imports) GetInvoice(idInvoice int) (*models.Invoice, error) {

	GetInvoiceById := i.dependencyManager.Get("IM__GetInvoiceById")

	requestContext := context.Background()
	requestContext = context.WithValue(requestContext, types.InvoiceId, idInvoice)

	invoice, err := GetInvoiceById(requestContext)

	if err != nil {
		return nil, err
	}

	return invoice.(*models.Invoice), err
}

func (i *Imports) CheckPaymentsOwnership(idCustomer int, paymentIds []int) error {

	CheckPaymentsOwnership := i.dependencyManager.Get("PM__CheckPaymentsOwnership")

	requestContext := context.Background()
	requestContext = context.WithValue(requestContext, types.CustomerId, idCustomer)
	requestContext = context.WithValue(requestContext, types.PaymentIds, paymentIds)

	_, err := CheckPaymentsOwnership(requestContext)

	if err != nil {
		return err
	}
	return nil
}

func (i *Imports) GetPayment(idPayment int) (*models.Payment, error) {

	GetPaymentById := i.dependencyManager.Get("PM__GetPaymentById")

	requestContext := context.Background()
	requestContext = context.WithValue(requestContext, types.PaymentId, idPayment)

	paymentRecord, err := GetPaymentById(requestContext)

	if err != nil {
		return nil, err
	}

	return paymentRecord.(*models.Payment), err

}

func (i *Imports) UpdatePayment(transaction *xorm.Session, payment *models.Payment) error {

	UpdatePayment := i.dependencyManager.Get("PM__UpdatePayment")

	requestContext := context.Background()
	requestContext = context.WithValue(requestContext, types.Transaction, transaction)
	requestContext = context.WithValue(requestContext, types.Payment, payment)

	_, err := UpdatePayment(requestContext)

	if err != nil {
		return err
	}

	return nil

}

func (i *Imports) UpdateInvoice(transaction *xorm.Session, invoice *models.Invoice) error {

	UpdateInvoice := i.dependencyManager.Get("IM__UpdateInvoice")

	requestContext := context.Background()
	requestContext = context.WithValue(requestContext, types.Transaction, transaction)
	requestContext = context.WithValue(requestContext, types.Invoice, invoice)

	_, err := UpdateInvoice(requestContext)

	if err != nil {
		return err
	}

	return nil
}
