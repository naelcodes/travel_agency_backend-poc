package service

import (
	"fmt"

	"neema.co.za/rest/utils/domains"
	"neema.co.za/rest/utils/logger"
	"neema.co.za/rest/utils/models"
	"neema.co.za/rest/utils/payloads"
	"neema.co.za/rest/utils/types"
)

func (s *Service) GetAllCustomerService__(queryParams *types.GetQueryParams) (*types.GetAllDTO[[]*models.Customer], error) {
	logger.Info("Getting all customers")
	return s.Repository.GetAll(queryParams)
}

func (s *Service) GetCustomerService(id int) (*models.Customer, error) {
	logger.Info("Getting customer")
	return s.Repository.GetById(id)
}

func (s *Service) CreateCustomerService(payload payloads.CreateCustomerPayload) (*models.Customer, error) {
	logger.Info(fmt.Sprintf("Creating customer: %v", payload))

	domain := domains.NewCustomerDomain(&payload.Customer)
	domain.SetDefaults()

	return s.Repository.Save(domain.GetCustomer())
}
