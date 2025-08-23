package service

import (
	"context"

	"neema.co.za/rest/utils/models"
	"neema.co.za/rest/utils/types"
	"xorm.io/xorm"
)

type Exports struct {
	InternalService *Service
}

func (e *Exports) IM__GetInvoiceById(context context.Context) (any, error) {

	idInvoice := context.Value(types.InvoiceId).(int)

	invoice, err := e.InternalService.GetById(idInvoice, nil, false)
	if err != nil {
		return nil, err
	}
	return invoice, nil
}

func (e *Exports) IM__UpdateInvoice(context context.Context) (any, error) {

	transaction := context.Value(types.Transaction).(*xorm.Session)
	invoice := context.Value(types.Invoice).(*models.Invoice)

	err := e.InternalService.Update(transaction, invoice)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
