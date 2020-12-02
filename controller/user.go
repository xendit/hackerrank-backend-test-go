package controller

import (
	"github.com/labstack/echo/v4"
)

type IUserService interface {
	// Task:
	// Fill this contract
	// e.g: Get(ctx context.Context, id string) (u entities.User, err error)
}

type userHandler struct {
	svc IUserService
}

func (h *userHandler) GetByID(c echo.Context) (err error) {
	panic("TODO")
}

// InitUserHandler will initiate the user handler
func InitUserHandler(e *echo.Echo, userSvc IUserService) {
	h := &userHandler{}
	e.GET("/users/:id", h.GetByID)
}
