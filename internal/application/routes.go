package application

import "github.com/labstack/echo/v4"

func registerRoutes(app *application, e *echo.Echo) {
	apiv1 := e.Group("/api/v1")

	library := apiv1.Group("/library")
	{
		library.GET("", app.libraryHandler.Get)
	}
}
