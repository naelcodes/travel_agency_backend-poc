//go:build wireinject

package customer

import (
	"github.com/google/wire"

	"neema.co.za/rest/modules/customer/internal/api"
	"neema.co.za/rest/modules/customer/internal/repository"
	"neema.co.za/rest/modules/customer/internal/service"
	"neema.co.za/rest/utils/app"
	"neema.co.za/rest/utils/database"
	"neema.co.za/rest/utils/managers"
)

// New api handler
func BuildApi(dependencyManager *managers.DependencyManager) *api.Api {
	panic(wire.Build(
		database.GetDatabase,
		app.NewFiberApp,
		repository.NewRepository,
		service.NewService,
		api.NewApi,
	))
}
