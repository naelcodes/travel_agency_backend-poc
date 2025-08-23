package api

import (
	"github.com/gofiber/fiber/v2"
	svc "neema.co.za/rest/modules/booking/internal/service"
)

type Module Api
type Api struct {
	*svc.Service
	*fiber.App
	*svc.Exports
}

func NewApi(service *svc.Service, app *fiber.App) *Api {
	return &Api{service, app, &svc.Exports{InternalService: service}}
}
