package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"neema.co.za/rest/utils/logger"
	"neema.co.za/rest/utils/types"
)

func (api *Api) GetAllTravelItemsHandler(c *fiber.Ctx) error {

	queryParams := c.Locals("queryParams").(*types.GetQueryParams)

	getAllTravelItemDTO, err := api.Service.GetAllTravelItemsService(queryParams)

	if err != nil {
		logger.Error(fmt.Sprintf("Error getting all travel items DTO: %v", err))
		return err
	}

	return c.Status(fiber.StatusOK).JSON(getAllTravelItemDTO)

}
