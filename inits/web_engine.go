package inits

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"nyacap-mini/middlewares"
	"nyacap-mini/routers"
)

func WebEngine() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:        true,
		LogStatus:     true,
		LogValuesFunc: middlewares.LogValues,
	}))

	routers.R(e)

	return e
}
