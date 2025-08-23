package errors

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

type CustomError struct {
	Type      string
	Err       error
	Operation string
}

func (c *CustomError) Error() string {
	return c.Err.Error()
}

func ServiceError(err error, operation string) *CustomError {
	return &CustomError{
		Err:       err,
		Type:      "ServiceError",
		Operation: operation,
	}
}

func DomainError(err error) *CustomError {
	return &CustomError{
		Err:  err,
		Type: "Domain Validation Error",
	}
}

func RepositoryError(err error) *CustomError {
	return &CustomError{
		Err:  err,
		Type: "RepositoryError",
	}
}

func NotFoundError(err error) *CustomError {
	return &CustomError{
		Err:  err,
		Type: "NotFoundError",
	}
}

func ValidationError(err error) *CustomError {
	return &CustomError{
		Err:  err,
		Type: "ValidationError",
	}
}

func UnknownError(err error) *CustomError {
	return &CustomError{
		Err:  err,
		Type: "UnknownError",
	}
}

func NewServerError(err error) *CustomError {
	return &CustomError{
		Err:  err,
		Type: "ServerError",
	}
}

type GlobalErrorResponse struct {
	Success   bool   `json:"success"`
	ErrorType string `json:"errorType,omitempty"`
	Error     string `json:"error,omitempty"`
	Operation string `json:"operation,omitempty"`
}

func ErrorResponse(code int, customError *CustomError, ctx *fiber.Ctx) error {
	return ctx.Status(code).JSON(GlobalErrorResponse{
		Success:   false,
		ErrorType: customError.Type,
		Error:     customError.Err.Error(),
		Operation: customError.Operation,
	})

}

func GlobalErrorHandler(ctx *fiber.Ctx, err error) error {

	var fiberError *fiber.Error
	if customError, ok := err.(*CustomError); ok {

		switch customError.Type {
		case "ValidationError":
			return ErrorResponse(fiber.StatusBadRequest, customError, ctx)
		case "RepositoryError":
			return ErrorResponse(fiber.StatusInternalServerError, customError, ctx)
		case "NotFoundError":
			return ErrorResponse(fiber.StatusBadRequest, customError, ctx)
		case "DomainError":
			return ErrorResponse(fiber.StatusBadRequest, customError, ctx)
		case "ServiceError":
			return ErrorResponse(fiber.StatusInternalServerError, customError, ctx)
		case "ServerError":
			return ErrorResponse(fiber.StatusInternalServerError, customError, ctx)
		case "UnknownError":
			return ErrorResponse(fiber.StatusInternalServerError, customError, ctx)
		}

	} else if errors.As(err, &fiberError) {
		return ctx.Status(fiberError.Code).JSON(GlobalErrorResponse{
			Success: false,
			Error:   fiberError.Error(),
		})
	} else {
		return ctx.Status(fiber.StatusInternalServerError).JSON(GlobalErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return nil

}
