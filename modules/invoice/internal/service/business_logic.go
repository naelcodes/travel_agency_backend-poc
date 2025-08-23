package service

import (
	"fmt"

	"neema.co.za/rest/utils/domains"
	CustomErrors "neema.co.za/rest/utils/errors"
	"neema.co.za/rest/utils/logger"
	"neema.co.za/rest/utils/managers"
	"neema.co.za/rest/utils/models"
	"neema.co.za/rest/utils/payloads"
	"neema.co.za/rest/utils/types"
)

func (s *Service) GetAllInvoiceService(queryParams *types.GetQueryParams) (*types.GetAllDTO[any], error) {
	logger.Info("Getting all invoices")
	return s.Repository.GetAll(queryParams)

}

func (s *Service) GetInvoiceService(id int, queryParams *types.GetQueryParams) (any, error) {
	logger.Info("Getting invoice")
	return s.Repository.GetById(id, queryParams, false)
}

func (s *Service) CreateInvoiceService(payload payloads.CreateInvoicePayload) (*models.Invoice, error) {
	logger.Info("Creating invoice")

	travelItems, err := s.Imports.GetTravelItems(payload.TravelItemIds)

	if err != nil {
		logger.Error(fmt.Sprintf("Error getting travel items: %v", err))
		return nil, err
	}

	invoiceAmount := float64(0)

	for _, travelItem := range travelItems {
		//check if it has been invoiced
		if travelItem.IdInvoice != 0 {
			logger.Error(fmt.Sprintf("Travel item %v has already been invoiced", travelItem.Id))
			return nil, CustomErrors.ValidationError(fmt.Errorf("travel item %v has already been invoiced", travelItem.Id))
		}
		invoiceAmount += travelItem.TotalPrice
	}

	invoiceDomain := domains.NewInvoiceDomain(&payload.Invoice)
	invoiceDomain.GetInvoice().Amount = invoiceAmount
	invoiceDomain.SetDefaults()

	err = invoiceDomain.Validate()

	if err != nil {
		logger.Error(fmt.Sprintf("Error validating invoice: %v", err))
		return nil, err
	}

	TransactionManager := managers.NewTransactionManager(s.Engine)
	err = TransactionManager.Begin()

	if err != nil {
		return nil, CustomErrors.UnknownError(err)
	}

	logger.Info(fmt.Sprintf("Invoice: %v", invoiceDomain.GetInvoice()))
	newInvoiceRecord, err := s.Repository.Save(TransactionManager.GetTransaction(), invoiceDomain.GetInvoice())

	if err != nil {
		TransactionManager.Rollback()
		return nil, err
	}

	err = s.Imports.LinkInvoiceToTravelItems(TransactionManager.GetTransaction(), newInvoiceRecord.Id, payload.TravelItemIds)

	if err != nil {
		TransactionManager.Rollback()
		return nil, err
	}

	TransactionManager.Commit()

	return newInvoiceRecord, nil
}
