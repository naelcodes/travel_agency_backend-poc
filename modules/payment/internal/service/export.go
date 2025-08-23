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

func (e *Exports) PM__CheckPaymentsOwnership(context context.Context) (any, error) {

	customerId := context.Value(types.CustomerId).(int)
	paymentIds := context.Value(types.PaymentIds).([]int)

	err := e.InternalService.CheckPaymentsOwnership(customerId, paymentIds)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (e *Exports) PM__GetPaymentById(context context.Context) (any, error) {

	paymentId := context.Value(types.PaymentId).(int)

	paymentRecord, err := e.InternalService.GetById(paymentId, nil)

	if err != nil {
		return nil, err
	}
	return paymentRecord, nil
}

func (e *Exports) PM__UpdatePayment(context context.Context) (any, error) {

	transaction := context.Value(types.Transaction).(*xorm.Session)
	paymentRecord := context.Value(types.Payment).(*models.Payment)

	err := e.InternalService.Update(transaction, paymentRecord)

	if err != nil {
		return nil, err
	}
	return nil, nil
}
