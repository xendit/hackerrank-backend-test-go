package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/xendit/hackerrank-backend-test-go/controller"
	"github.com/xendit/hackerrank-backend-test-go/repositories"
	"github.com/xendit/hackerrank-backend-test-go/services"

	// Imported for migrations and database driver
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

const (
	defaultAppPort = ":8000"
)

func main() {
	// open database
	dbConn, err := sql.Open("sqlite3", repositories.SqliteDBDsn)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// Run migrations
	migrator, err := repositories.RunMigrationSQLite(dbConn, "./repositories/migrations")
	if err != nil {
		log.Fatal(err)
	}
	_, err = migrator.Up()
	if err != nil {
		log.Fatal(err)
	}

	usrRepo := repositories.NewUser()
	usrSvc := services.NewUser(usrRepo)

	e := echo.New()
	controller.InitHealthCheckHandler(e)
	controller.InitUserHandler(e, usrSvc)
	logrus.Fatal(e.Start(defaultAppPort))
}
