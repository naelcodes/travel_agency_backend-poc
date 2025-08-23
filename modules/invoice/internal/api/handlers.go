package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	CustomErrors "neema.co.za/rest/utils/errors"
	"neema.co.za/rest/utils/logger"
	"neema.co.za/rest/utils/payloads"
	"neema.co.za/rest/utils/types"
)

func (api *Api) GetAllInvoiceHandler(c *fiber.Ctx) error {
	queryParams := c.Locals("queryParams").(*types.GetQueryParams)

	invoicesDTO, err := api.Service.GetAllInvoiceService(queryParams)

	if err != nil {
		logger.Error(fmt.Sprintf("Error getting all invoices DTO: %v", err))
		return err
	}

	logger.Info(fmt.Sprintf("All invoices DTO: %v", invoicesDTO))

	return c.Status(fiber.StatusOK).JSON(invoicesDTO)
}

func (api *Api) GetInvoiceHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	if err != nil {
		logger.Error(fmt.Sprintf("Error parsing id: %v", err))
		return CustomErrors.ServiceError(err, "parsing id")
	}

	logger.Info(fmt.Sprintf("params Id: %v", id))

	queryParams := c.Locals("queryParams").(*types.GetQueryParams)

	invoiceDTO, err := api.Service.GetInvoiceService(id, queryParams)

	if err != nil {
		logger.Error(fmt.Sprintf("Error getting invoice DTO: %v", err))
		return err
	}

	logger.Info(fmt.Sprintf("Invoice DTO: %v", invoiceDTO))

	return c.Status(fiber.StatusOK).JSON(invoiceDTO)
}

func (api *Api) CreateInvoiceHandler(c *fiber.Ctx) error {
	createInvoicePayload := c.Locals("payload").(*payloads.CreateInvoicePayload)
	logger.Info(fmt.Sprintf("CreateInvoiceDTO: %v", createInvoicePayload))

	newInvoiceRecord, err := api.Service.CreateInvoiceService(*createInvoicePayload)

	if err != nil {
		logger.Error(fmt.Sprintf("Error creating invoice Record: %v", err))
		return err
	}

	logger.Info(fmt.Sprintf("NewInvoiceDTO: %v", newInvoiceRecord))
	return c.Status(fiber.StatusCreated).JSON(newInvoiceRecord)
}
