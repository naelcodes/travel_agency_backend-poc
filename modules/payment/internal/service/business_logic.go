package service

import (
	"fmt"

	"neema.co.za/rest/utils/domains"
	"neema.co.za/rest/utils/managers"
	"neema.co.za/rest/utils/models"
	"neema.co.za/rest/utils/payloads"

	CustomErrors "neema.co.za/rest/utils/errors"
	"neema.co.za/rest/utils/logger"
	"neema.co.za/rest/utils/types"
)

func (s *Service) GetAllPaymentsService(queryParams *types.GetQueryParams) (*types.GetAllDTO[any], error) {
	logger.Info("Getting all payments")
	return s.Repository.GetAll(queryParams)
}

func (s *Service) GetPaymentService(id int, queryParams *types.GetQueryParams) (any, error) {
	logger.Info("Getting payment")
	return s.Repository.GetById(id, queryParams)
}

func (s *Service) CreatePaymentService(payload payloads.CreatePaymentPayload) (*models.Payment, error) {
	logger.Info("Creating payment")

	paymentDomain := domains.NewPaymentDomain(&payload.Payment)
	paymentDomain.SetDefaults()

	err := paymentDomain.Validate()

	if err != nil {
		logger.Error(fmt.Sprintf("payment validation error: %v", err))
		return nil, err
	}

	TransactionManager := managers.NewTransactionManager(s.Engine)
	err = TransactionManager.Begin()

	if err != nil {
		return nil, CustomErrors.UnknownError(err)
	}

	logger.Info(fmt.Sprintf("Payment: %v", paymentDomain.GetPayment()))

	paymentRecord, err := s.Repository.Save(TransactionManager.GetTransaction(), paymentDomain.GetPayment())

	if err != nil {
		TransactionManager.Rollback()
		return nil, err
	}

	TransactionManager.Commit()

	return paymentRecord, nil

}
