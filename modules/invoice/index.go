package invoice

import (
	. "neema.co.za/rest/modules/invoice/internal/api"
	"neema.co.za/rest/utils/managers"
	"neema.co.za/rest/utils/middlewares"
	"neema.co.za/rest/utils/payloads"
)

func GetModule(dependencyManager *managers.DependencyManager) *Module {
	api := BuildApi(dependencyManager)
	handleRoutes(api)
	module := Module(*api)                //Module is an alias of Api
	dependencyManager.Add(module.Exports) //add exportable functions to dependency manager
	return &module
}

func handleRoutes(api *Api) {
	api.Get("", api.GetAllInvoiceHandler)
	api.Get("/:id", api.GetInvoiceHandler)
	api.Post("", middlewares.PayloadValidator(new(payloads.CreateInvoicePayload)), api.CreateInvoiceHandler)

}
