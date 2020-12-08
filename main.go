package main

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/xendit/hackerrank-backend-test-go/controller"
	"github.com/xendit/hackerrank-backend-test-go/repositories"
	"github.com/xendit/hackerrank-backend-test-go/services"
)

const (
	defaultAppPort = ":8000"
)

func main() {
	usrRepo := repositories.NewUser()
	usrSvc := services.NewUser(usrRepo)

	e := echo.New()
	controller.InitHealthCheckHandler(e)
	controller.InitUserHandler(e, usrSvc)
	logrus.Fatal(e.Start(defaultAppPort))
}
