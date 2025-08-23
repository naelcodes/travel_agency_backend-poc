package middlewares

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	CustomErrors "neema.co.za/rest/utils/errors"
	"neema.co.za/rest/utils/logger"
	"neema.co.za/rest/utils/payloads"
	"neema.co.za/rest/utils/types"
)

func GetCors() fiber.Handler {
	config := cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE",
	}

	return cors.New(config)
}

func Recover() fiber.Handler {
	config := recover.Config{EnableStackTrace: true}
	return recover.New(config)
}

func QueryValidator() fiber.Handler {
	return func(c *fiber.Ctx) error {
		queryParams := new(types.GetQueryParams)
		err := c.QueryParser(queryParams)

		if err != nil {
			return CustomErrors.ServiceError(err, "Parsing query params")
		}

		if c.Method() == fiber.MethodGet {
			if queryParams.PageNumber != nil && queryParams.PageSize == nil {
				return CustomErrors.ValidationError(errors.New("page size should be provided with page number"))
			}

			if queryParams.PageSize != nil && queryParams.PageNumber == nil {
				return CustomErrors.ValidationError(errors.New("page number should be provided with page size"))
			}

			if queryParams.PageSize != nil && queryParams.PageNumber != nil {
				if *queryParams.PageSize <= 0 {
					return CustomErrors.ValidationError(errors.New("page size should be greater than 0"))
				}

				if *queryParams.PageNumber < 0 {
					return CustomErrors.ValidationError(errors.New("page number should be greater than or equal to 0"))
				}
			}
		}

		c.Locals("queryParams", queryParams)

		return c.Next()
	}
}

func PayloadValidator(payload types.PayloadValidator) fiber.Handler {
	return func(c *fiber.Ctx) error {

		logger.Info(fmt.Sprintf("Validating payload on path: %v/ %v", c.Method(), c.Path()))

		if c.Method() == fiber.MethodPost || c.Method() == fiber.MethodPut || c.Method() == fiber.MethodPatch {

			err := c.BodyParser(payload)

			if err != nil {
				return CustomErrors.ServiceError(err, "Parsing JSON payload")
			}

			if err = payload.Validate(); err != nil {
				return CustomErrors.ValidationError(err)
			}

			c.Locals("payload", payload)
		}
		return c.Next()
	}
}

func ImputationPayloadValidator() fiber.Handler {
	return func(c *fiber.Ctx) error {
		payload := []*payloads.ImputationPayload{}

		logger.Info(fmt.Sprintf("Validating payload list on path: %v/ %v", c.Method(), c.Path()))

		if c.Method() == fiber.MethodPost || c.Method() == fiber.MethodPut || c.Method() == fiber.MethodPatch {

			err := c.BodyParser(&payload)

			if err != nil {
				return CustomErrors.ServiceError(err, "Parsing JSON payload list ### ")
			}

			if len(payload) == 0 {
				return CustomErrors.ValidationError(errors.New("payload list should not be empty"))
			}

			for _, p := range payload {
				if err = p.Validate(); err != nil {
					return CustomErrors.ValidationError(err)
				}
			}

			c.Locals("payload", payload)
		}
		return c.Next()

	}
}
