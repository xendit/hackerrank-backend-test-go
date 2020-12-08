package e2e

import (
	"database/sql"
	"net/http"

	"github.com/stretchr/testify/suite"
	"github.com/xendit/hackerrank-backend-test-go/repositories"

	// This is imported for migrations
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type Suite struct {
	suite.Suite
	DBConn                  *sql.DB
	Client                  *http.Client
	Migration               *repositories.Migration
	DBDsn                   string
	MigrationLocationFolder string
}

// SetupSuite setup at the beginning of test
func (s *Suite) SetupSuite() {
	var err error
	s.DBConn, err = sql.Open("postgres", s.DBDsn)
	s.Require().NoError(err)
	err = s.DBConn.Ping()
	s.Require().NoError(err)
	s.Migration, err = repositories.RunMigrationPostgres(s.DBConn, s.MigrationLocationFolder)
	s.Require().NoError(err)

	s.Client = http.DefaultClient
}

// TearDownSuite teardown at the end of test
func (s *Suite) TearDownSuite() {
	err := s.DBConn.Close()
	s.Require().NoError(err)
}
