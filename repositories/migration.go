package repositories

import (
	"database/sql"
	"strings"

	migrate "github.com/golang-migrate/migrate/v4"
	_mysql "github.com/golang-migrate/migrate/v4/database/mysql"
	_postgres "github.com/golang-migrate/migrate/v4/database/postgres"
)

type Migration struct {
	Migrate *migrate.Migrate
}

func (m *Migration) Up() (bool, error) {
	err := m.Migrate.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			return true, nil
		}
		return false, err
	}
	return true, nil
}

func (m *Migration) Down() (bool, error) {
	err := m.Migrate.Down()
	if err != nil {
		return false, err
	}
	return true, err
}

func RunMigrationPostgres(dbConn *sql.DB, migrationsFolderLocation string) (*Migration, error) {
	dataPath := []string{}
	dataPath = append(dataPath, "file://")
	dataPath = append(dataPath, migrationsFolderLocation)

	pathToMigrate := strings.Join(dataPath, "")

	driver, err := _postgres.WithInstance(dbConn, &_postgres.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(pathToMigrate, "postgres", driver)
	if err != nil {
		return nil, err
	}
	return &Migration{Migrate: m}, nil
}

func RunMigrationMySQL(dbConn *sql.DB, migrationsFolderLocation string) (*Migration, error) {
	dataPath := []string{}
	dataPath = append(dataPath, "file://")
	dataPath = append(dataPath, migrationsFolderLocation)

	pathToMigrate := strings.Join(dataPath, "")

	driver, err := _mysql.WithInstance(dbConn, &_mysql.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(pathToMigrate, "mysql", driver)
	if err != nil {
		return nil, err
	}

	return &Migration{Migrate: m}, nil
}
