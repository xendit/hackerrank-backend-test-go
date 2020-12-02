package controller

import (
	"github.com/labstack/echo/v4"
)

type healthCheckHandler struct {
}

func (h *healthCheckHandler) HealtcheckLiveness(c echo.Context) (err error) {
	panic("TODO")
}

// InitHealthCheckHandler will initiate the healtcheck handler
func InitHealthCheckHandler(e *echo.Echo) {
	h := &healthCheckHandler{}
	e.GET("/healthcheck/liveness", h.HealtcheckLiveness)
	e.GET("/healthcheck/readiness", nil) // TODO: add readiness handler
}
