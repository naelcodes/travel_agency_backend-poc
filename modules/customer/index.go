package customer

import (
	. "neema.co.za/rest/modules/customer/internal/api"
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
	api.Get("", api.GetAllCustomerHandler)
	api.Get("/:id", api.GetCustomerHandler)
	api.Post("", middlewares.PayloadValidator(new(payloads.CreateCustomerPayload)), api.CreateCustomerHandler)

}
