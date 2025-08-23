package imputation

import (
	. "neema.co.za/rest/modules/imputation/internal/api"
	"neema.co.za/rest/utils/managers"
	"neema.co.za/rest/utils/middlewares"
)

func GetModule(dependencyManager *managers.DependencyManager) *Module {
	api := BuildApi(dependencyManager)
	handleRoutes(api)
	module := Module(*api)                //Module is an alias of Api
	dependencyManager.Add(module.Exports) //add exportable functions to dependency manager
	return &module
}

func handleRoutes(api *Api) {
	api.Get("/:id/imputations", api.GetImputationsHandler)
	api.Post("/:id/imputations", middlewares.ImputationPayloadValidator(), api.ApplyImputationsHandler)

}
