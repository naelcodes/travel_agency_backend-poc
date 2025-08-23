package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	CustomErrors "neema.co.za/rest/utils/errors"
	"neema.co.za/rest/utils/logger"
	"neema.co.za/rest/utils/payloads"
)

func (api *Api) GetImputationsHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	if err != nil {
		logger.Error(fmt.Sprintf("Error parsing id: %v", err))
		return CustomErrors.ServiceError(err, "parsing id")
	}

	logger.Info(fmt.Sprintf("params Id: %v", id))

	invoiceImputationRecords, err := api.Service.GetImputationsService(id)

	if err != nil {
		logger.Error(fmt.Sprintf("Error getting invoice imputations: %v", err))
		return err
	}

	return c.Status(fiber.StatusOK).JSON(invoiceImputationRecords)

}

func (api *Api) ApplyImputationsHandler(c *fiber.Ctx) error {
	idInvoice, err := c.ParamsInt("id")

	if err != nil {
		logger.Error(fmt.Sprintf("Error parsing id: %v", err))
		return CustomErrors.ServiceError(err, "parsing id")
	}

	logger.Info(fmt.Sprintf("params Id: %v", idInvoice))

	payload := c.Locals("payload").([]*payloads.ImputationPayload)

	insertedCount, UpdatedCount, deletedCount, err := api.Service.ApplyImputationsService(idInvoice, payload)

	if err != nil {
		logger.Error(fmt.Sprintf("Error applying imputations: %v", err))
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Inserted Imputation Count": insertedCount,
		"Updated Imputation Count":  UpdatedCount,
		"Deleted Imputation Count":  deletedCount,
	})

}
