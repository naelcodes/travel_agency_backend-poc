package service

import (
	"neema.co.za/rest/utils/logger"
	"neema.co.za/rest/utils/models"
	"neema.co.za/rest/utils/types"
)

func (s *Service) GetAllTravelItemsService(queryParams *types.GetQueryParams) (*types.GetAllDTO[[]*models.TravelItem], error) {
	logger.Info("Getting all travel items")
	return s.Repository.GetAll(queryParams)
}

// func (s *Service) AddInvoiceToTravelItemService(transaction *xorm.Session, invoiceId int, travelItemIds []int) error {
// 	logger.Info("Adding invoice to travel item")
// 	return s.Repository.AddInvoiceToTravelItem(transaction, invoiceId, travelItemIds)
// }
