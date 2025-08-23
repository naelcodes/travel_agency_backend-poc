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

func (i *Imports) LinkInvoiceToTravelItems(transaction *xorm.Session, IdInvoice int, travelItemIds []int) error {

	InvoiceTravelItems := i.dependencyManager.Get("BM__InvoiceTravelItems")

	requestContext := context.Background()
	requestContext = context.WithValue(requestContext, types.InvoiceId, IdInvoice)
	requestContext = context.WithValue(requestContext, types.TravelItemIds, travelItemIds)
	requestContext = context.WithValue(requestContext, types.Transaction, transaction)

	_, err := InvoiceTravelItems(requestContext)

	if err != nil {
		return err
	}

	return nil
}

func (i *Imports) GetTravelItems(travelItemIds []int) ([]*models.TravelItem, error) {

	GetTravelItemsByIds := i.dependencyManager.Get("BM__GetTravelItemsByIds")

	requestContext := context.Background()
	requestContext = context.WithValue(requestContext, types.TravelItemIds, travelItemIds)
	travelItems, err := GetTravelItemsByIds(requestContext)

	if err != nil {
		return nil, err
	}
	return travelItems.([]*models.TravelItem), nil
}
