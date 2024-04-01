package routers

import (
	"github.com/labstack/echo/v4"
	"nyacap-mini/handlers/public"
)

func Public(rg *echo.Group) {
	rg.GET("/healthcheck", public.HealthCheck)
}
