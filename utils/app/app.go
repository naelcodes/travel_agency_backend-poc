package app

import (
	"github.com/gofiber/fiber/v2"
	"neema.co.za/rest/utils/errors"
	"neema.co.za/rest/utils/middlewares"
)

// type RouterCreator func(prefix string, handlers ...fiber.Handler) fiber.Router
// type Router fiber.Router
type App struct {
	*fiber.App
}

var app *App

func init() {
	app = &App{NewFiberApp()}
}

func NewFiberApp() *fiber.App {
	config := fiber.Config{
		ErrorHandler: errors.GlobalErrorHandler,
	}

	return fiber.New(config)
}

//func NewRouter() RouterCreator {
//	return app.Group
//}

func Initialise() *App {

	app.Use(middlewares.Recover())
	app.Use(middlewares.GetCors())
	app.Use(middlewares.QueryValidator())

	return app
}
