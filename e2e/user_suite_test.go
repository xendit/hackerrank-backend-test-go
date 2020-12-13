package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/xendit/hackerrank-backend-test-go/repositories"
)

const (
	DefaultHost = "http://localhost:8000"
)

var (
	RequiredHeaders = http.Header{
		"client-version": []string{"v0.0.0"},
		"team-name":      []string{"test"},
		"service-name":   []string{"local"},
		"Content-Type":   []string{"application/json"},
	}
)

type userTestSuite struct {
	Suite
	Host string
}

func TestUserSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip the Test Suite for Payment Repository")
	}

	host := os.Getenv("TEST_HOST")
	if host == "" {
		host = DefaultHost
	}

	userSuite := &userTestSuite{
		Suite: Suite{
			DBDsn:                   fmt.Sprintf("../%s", repositories.SqliteDBDsn),
			MigrationLocationFolder: "../repositories/migrations",
		},
		Host: host,
	}

	suite.Run(t, userSuite)
}

func (s userTestSuite) BeforeTest(_, _ string) {
	ok, err := s.Migration.Up()
	s.Require().NoError(err)
	s.Require().True(ok)
}

func (s userTestSuite) AfterTest(_, _ string) {
	ok, err := s.Migration.Down()
	s.Require().NoError(err)
	s.Require().True(ok)
}

func (s userTestSuite) seedFetchUser() {
	uri := fmt.Sprintf("%s/api/users", s.Host)
	users := []map[string]interface{}{
		{
			"firstName": "First",
			"lastName":  "User",
			"address":   "Indonesia",
			"isActive":  true,
		},
		{
			"firstName": "Second",
			"lastName":  "User",
			"address":   "Indonesia",
			"isActive":  true,
		},
		{
			"firstName": "Third",
			"lastName":  "User",
			"address":   "Indonesia",
			"isActive":  true,
		},
	}

	for _, item := range users {
		jbyt, err := json.Marshal(item)
		s.Require().NoError(err)

		req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(jbyt))
		s.Require().NoError(err)

		req.Header = RequiredHeaders
		resp, err := s.Client.Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusCreated, resp.StatusCode)
		respByte, err := ioutil.ReadAll(resp.Body)
		s.Require().NoError(err)
		defer resp.Body.Close()
		respMap := map[string]interface{}{}
		err = json.Unmarshal(respByte, &respMap)
		s.Require().NoError(err)
		s.Require().NotNil(respMap)
		s.Require().NotNil(respMap["id"])
		s.Require().NotEmpty(respMap["id"])
		s.Require().Equal(item["firstName"], respMap["firstName"])
		s.Require().Equal(item["lastName"], respMap["lastName"])
		s.Require().Equal(item["address"], respMap["address"])
		s.Require().Equal(item["isActive"], respMap["isActive"])
	}
}

func (s userTestSuite) TestFetchUser() {
	s.T().Log("Seeding the user data")
	s.seedFetchUser()

	firstUrl := fmt.Sprintf("%s/api/users?limit=2&offset=0", s.Host)
	req, err := http.NewRequest(http.MethodGet, firstUrl, nil)
	s.Require().NoError(err)
	req.Header = RequiredHeaders

	resp, err := s.Client.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	respByte, err := ioutil.ReadAll(resp.Body)
	s.Require().NoError(err)
	defer resp.Body.Close()
	respMap := []map[string]interface{}{}
	err = json.Unmarshal(respByte, &respMap)
	s.Require().NoError(err)
	s.Require().NotNil(respMap)
	s.Require().Len(respMap, 2)
	s.Require().Equal("Third", respMap[0]["firstName"])
	s.Require().Equal("Second", respMap[1]["firstName"])

	secondUrl := fmt.Sprintf("%s/api/users?limit=2&offset=2", s.Host)
	req, err = http.NewRequest(http.MethodGet, secondUrl, nil)
	s.Require().NoError(err)
	req.Header = RequiredHeaders

	resp, err = s.Client.Do(req)
	s.Require().NoError(err)

	respByte, err = ioutil.ReadAll(resp.Body)
	s.Require().NoError(err)
	defer resp.Body.Close()

	respMap = []map[string]interface{}{}
	err = json.Unmarshal(respByte, &respMap)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	s.Require().NotNil(respMap)
	s.Require().Len(respMap, 1)
	s.Require().Equal("First", respMap[0]["firstName"])
}

func (s userTestSuite) TestCreateUser() {
	uri := fmt.Sprintf("%s/api/users", s.Host)

	s.T().Run("Return Status Created", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"firstName": "John",
			"lastName":  "Doe",
			"address":   "Singapore",
			"isActive":  true,
		}
		jbyt, err := json.Marshal(reqBody)
		s.Require().NoError(err)

		req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(jbyt))
		s.Require().NoError(err)

		req.Header = RequiredHeaders
		resp, err := s.Client.Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusCreated, resp.StatusCode)
		respByte, err := ioutil.ReadAll(resp.Body)
		s.Require().NoError(err)
		defer resp.Body.Close()

		respMap := map[string]interface{}{}
		err = json.Unmarshal(respByte, &respMap)
		s.Require().NoError(err)
		s.Require().NotNil(respMap)
		s.assertUserValue(reqBody, respMap)
	})

	s.T().Run("Return 400 for incorrect field type", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"firstName": "John",
			"lastName":  "Doe",
			"address":   "Singapore",
			"isActive":  "potato",
		}
		jbyt, err := json.Marshal(reqBody)
		s.Require().NoError(err)

		req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(jbyt))
		s.Require().NoError(err)

		req.Header = RequiredHeaders
		resp, err := s.Client.Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
		respByte, err := ioutil.ReadAll(resp.Body)
		s.Require().NoError(err)
		defer resp.Body.Close()

		respMap := map[string]interface{}{}
		err = json.Unmarshal(respByte, &respMap)
		s.Require().NoError(err)
		s.Require().NotNil(respMap)
		s.Require().Equal("API_VALIDATION_ERROR", respMap["error_code"])
	})

	s.T().Run("Return 400 for missing required field", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"lastName": "Doe",
			"address":  "Singapore",
			"isActive": true,
		}
		jbyt, err := json.Marshal(reqBody)
		s.Require().NoError(err)

		req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(jbyt))
		s.Require().NoError(err)

		req.Header = RequiredHeaders
		resp, err := s.Client.Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
		respByte, err := ioutil.ReadAll(resp.Body)
		s.Require().NoError(err)
		defer resp.Body.Close()

		respMap := map[string]interface{}{}
		err = json.Unmarshal(respByte, &respMap)
		s.Require().NoError(err)
		s.Require().NotNil(respMap)
		s.Require().Equal("API_VALIDATION_ERROR", respMap["error_code"])
	})
}

func (s userTestSuite) TestGetUserByID() {
	// Seed user data
	uri := fmt.Sprintf("%s/api/users", s.Host)
	reqBody := map[string]interface{}{
		"firstName": "John",
		"lastName":  "Doe",
		"address":   "Singapore",
		"isActive":  true,
	}
	jbyt, err := json.Marshal(reqBody)
	s.Require().NoError(err)

	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(jbyt))
	s.Require().NoError(err)

	req.Header = RequiredHeaders
	resp, err := s.Client.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusCreated, resp.StatusCode)
	respByte, err := ioutil.ReadAll(resp.Body)
	s.Require().NoError(err)
	defer resp.Body.Close()

	respMap := map[string]interface{}{}
	err = json.Unmarshal(respByte, &respMap)
	s.Require().NoError(err)
	s.assertUserValue(reqBody, respMap)
	// new request for GET
	getUserByIDURI := fmt.Sprintf("%s/api/users/%v", s.Host, respMap["id"])
	req, err = http.NewRequest(http.MethodGet, getUserByIDURI, nil)
	s.Require().NoError(err)

	req.Header = RequiredHeaders
	resp, err = s.Client.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	respByte, err = ioutil.ReadAll(resp.Body)
	s.Require().NoError(err)
	defer resp.Body.Close()

	respMap = map[string]interface{}{}
	err = json.Unmarshal(respByte, &respMap)
	s.Require().NoError(err)
	s.assertUserValue(reqBody, respMap)

}

func (s userTestSuite) TestUpdateUserUsingPUT() {
	s.T().Log("Seeding the user data")
	uri := fmt.Sprintf("%s/api/users", s.Host)
	reqBody := map[string]interface{}{
		"firstName": "John",
		"lastName":  "Doe",
		"address":   "Singapore",
		"isActive":  true,
	}
	jbyt, err := json.Marshal(reqBody)
	s.Require().NoError(err)
	s.T().Log("Start the HTTP call")
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(jbyt))
	s.Require().NoError(err)

	req.Header = RequiredHeaders
	resp, err := s.Client.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusCreated, resp.StatusCode)
	respByte, err := ioutil.ReadAll(resp.Body)
	s.Require().NoError(err)
	defer resp.Body.Close()

	respMap := map[string]interface{}{}
	err = json.Unmarshal(respByte, &respMap)
	s.Require().NoError(err)
	s.assertUserValue(reqBody, respMap)

	s.T().Log("Ensure the inserted user data")
	// new request for GET
	getUserByIDURI := fmt.Sprintf("%s/api/users/%v", s.Host, respMap["id"])
	req, err = http.NewRequest(http.MethodGet, getUserByIDURI, nil)
	s.Require().NoError(err)

	req.Header = RequiredHeaders
	resp, err = s.Client.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	respByte, err = ioutil.ReadAll(resp.Body)
	s.Require().NoError(err)
	defer resp.Body.Close()

	respMap = map[string]interface{}{}
	err = json.Unmarshal(respByte, &respMap)
	s.Require().NoError(err)
	s.assertUserValue(reqBody, respMap)

	s.T().Log("Update the user data")
	updateUserByIDURI := fmt.Sprintf("%s/api/users/%v", s.Host, respMap["id"])
	reqBodyUpdate := map[string]interface{}{
		"firstName": "John",
		"lastName":  "Doe",
		"address":   "Jakarta",
		"isActive":  true,
	}
	jbyt, err = json.Marshal(reqBodyUpdate)
	s.Require().NoError(err)

	req, err = http.NewRequest(http.MethodPut, updateUserByIDURI, bytes.NewBuffer(jbyt))
	s.Require().NoError(err)

	req.Header = RequiredHeaders
	resp, err = s.Client.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	respByte, err = ioutil.ReadAll(resp.Body)
	s.Require().NoError(err)
	defer resp.Body.Close()

	updateRespMap := map[string]interface{}{}
	err = json.Unmarshal(respByte, &updateRespMap)
	s.Require().NoError(err)
	s.assertUserValue(reqBodyUpdate, updateRespMap)

}

func (s userTestSuite) TestDeleteUserByID() {
	// Seed user data
	uri := fmt.Sprintf("%s/api/users", s.Host)
	reqBody := map[string]interface{}{
		"firstName": "John",
		"lastName":  "Doe",
		"address":   "Singapore",
		"isActive":  true,
	}
	jbyt, err := json.Marshal(reqBody)
	s.Require().NoError(err)

	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(jbyt))
	s.Require().NoError(err)

	req.Header = RequiredHeaders
	resp, err := s.Client.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusCreated, resp.StatusCode)
	respByte, err := ioutil.ReadAll(resp.Body)
	s.Require().NoError(err)
	defer resp.Body.Close()

	respMap := map[string]interface{}{}
	err = json.Unmarshal(respByte, &respMap)
	s.Require().NoError(err)
	s.assertUserValue(reqBody, respMap)
	// new request for GET: assert the inserted User is exists
	getUserByIDURI := fmt.Sprintf("%s/api/users/%v", s.Host, respMap["id"])
	req, err = http.NewRequest(http.MethodGet, getUserByIDURI, nil)
	s.Require().NoError(err)

	req.Header = RequiredHeaders
	resp, err = s.Client.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	respByte, err = ioutil.ReadAll(resp.Body)
	s.Require().NoError(err)
	defer resp.Body.Close()

	respMap = map[string]interface{}{}
	err = json.Unmarshal(respByte, &respMap)
	s.Require().NoError(err)
	s.assertUserValue(reqBody, respMap)

	// new request for DELETE: assert the inserted User is exists
	deleteUserByIDURI := fmt.Sprintf("%s/api/users/%v", s.Host, respMap["id"])
	req, err = http.NewRequest(http.MethodDelete, deleteUserByIDURI, nil)
	s.Require().NoError(err)

	req.Header = RequiredHeaders
	resp, err = s.Client.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNoContent, resp.StatusCode)
	closerDelete := ioutil.NopCloser(resp.Body)
	s.Require().NoError(err)
	defer closerDelete.Close()

	// new request for GET: assert the delete action is success User is exists
	getDeletedUserByIDURI := fmt.Sprintf("%s/api/users/%v", s.Host, respMap["id"])
	req, err = http.NewRequest(http.MethodGet, getDeletedUserByIDURI, nil)
	s.Require().NoError(err)

	req.Header = RequiredHeaders
	resp, err = s.Client.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
	closerGet := ioutil.NopCloser(resp.Body)
	s.Require().NoError(err)
	defer closerGet.Close()
}

func (s *userTestSuite) assertUserValue(expected map[string]interface{}, actual map[string]interface{}) {
	s.Require().NotNil(actual)
	s.Require().NotNil(actual["id"])
	s.Require().NotEmpty(actual["id"])
	s.Require().Equal(expected["firstName"], actual["firstName"])
	s.Require().Equal(expected["lastName"], actual["lastName"])
	s.Require().Equal(expected["address"], actual["address"])
	s.Require().Equal(expected["isActive"], actual["isActive"])
}
