package service

import (
	"context"

	"neema.co.za/rest/utils/types"
	"xorm.io/xorm"
)

type Exports struct {
	InternalService *Service
}

func (e *Exports) BM__InvoiceTravelItems(context context.Context) (any, error) {

	transaction := context.Value(types.Transaction).(*xorm.Session)
	invoiceId := context.Value(types.InvoiceId).(int)
	travelItemIds := context.Value(types.TravelItemIds).([]int)

	err := e.InternalService.InvoiceTravelItems(transaction, invoiceId, travelItemIds)

	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (e *Exports) BM__GetTravelItemsByIds(context context.Context) (any, error) {

	travelItemIds := context.Value(types.TravelItemIds).([]int)
	travelItems, err := e.InternalService.GetByIds(travelItemIds)

	if err != nil {
		return nil, err
	}
	return travelItems, nil

}
