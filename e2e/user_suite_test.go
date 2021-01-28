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
		t.Skip("Skipping Test Suite")
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
	s.Assert().NoError(err)
	s.Assert().True(ok)
}

func (s userTestSuite) AfterTest(_, _ string) {
	ok, err := s.Migration.Down()
	s.Assert().NoError(err)
	s.Assert().True(ok)
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
		s.Assert().NoError(err)

		req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(jbyt))
		s.Assert().NoError(err)

		req.Header = RequiredHeaders
		resp, err := s.Client.Do(req)
		s.Assert().NoError(err)
		s.Assert().Equal(http.StatusCreated, resp.StatusCode)
		respByte, err := ioutil.ReadAll(resp.Body)
		s.Assert().NoError(err)
		defer resp.Body.Close()
		respMap := map[string]interface{}{}
		err = json.Unmarshal(respByte, &respMap)
		s.Assert().NoError(err)
		s.Assert().NotNil(respMap)
		s.Assert().NotNil(respMap["id"])
		s.Assert().NotEmpty(respMap["id"])
		s.Assert().Equal(item["firstName"], respMap["firstName"])
		s.Assert().Equal(item["lastName"], respMap["lastName"])
		s.Assert().Equal(item["address"], respMap["address"])
		s.Assert().Equal(item["isActive"], respMap["isActive"])
	}
}

func (s userTestSuite) TestFetchUser() {
	s.seedFetchUser()

	firstUrl := fmt.Sprintf("%s/api/users?limit=2&offset=0", s.Host)
	req, err := http.NewRequest(http.MethodGet, firstUrl, nil)
	s.Assert().NoError(err)
	req.Header = RequiredHeaders

	resp, err := s.Client.Do(req)
	s.Assert().NoError(err)
	s.Assert().Equal(http.StatusOK, resp.StatusCode)
	respByte, err := ioutil.ReadAll(resp.Body)
	s.Assert().NoError(err)
	defer resp.Body.Close()
	respMap := []map[string]interface{}{}
	err = json.Unmarshal(respByte, &respMap)
	s.Assert().NoError(err)
	s.Assert().NotNil(respMap)
	s.Assert().Len(respMap, 2)
	s.Assert().Equal("Third", respMap[0]["firstName"])
	s.Assert().Equal("Second", respMap[1]["firstName"])

	secondUrl := fmt.Sprintf("%s/api/users?limit=2&offset=2", s.Host)
	req, err = http.NewRequest(http.MethodGet, secondUrl, nil)
	s.Assert().NoError(err)
	req.Header = RequiredHeaders

	resp, err = s.Client.Do(req)
	s.Assert().NoError(err)

	respByte, err = ioutil.ReadAll(resp.Body)
	s.Assert().NoError(err)
	defer resp.Body.Close()

	respMap = []map[string]interface{}{}
	err = json.Unmarshal(respByte, &respMap)
	s.Assert().NoError(err)
	s.Assert().Equal(http.StatusOK, resp.StatusCode)
	s.Assert().NotNil(respMap)
	s.Assert().Len(respMap, 1)
	s.Assert().Equal("First", respMap[0]["firstName"])
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
		s.Assert().NoError(err)

		req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(jbyt))
		s.Assert().NoError(err)

		req.Header = RequiredHeaders
		resp, err := s.Client.Do(req)
		s.Assert().NoError(err)
		s.Assert().Equal(http.StatusCreated, resp.StatusCode)
		respByte, err := ioutil.ReadAll(resp.Body)
		s.Assert().NoError(err)
		defer resp.Body.Close()

		respMap := map[string]interface{}{}
		err = json.Unmarshal(respByte, &respMap)
		s.Assert().NoError(err)
		s.Assert().NotNil(respMap)
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
		s.Assert().NoError(err)

		req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(jbyt))
		s.Assert().NoError(err)

		req.Header = RequiredHeaders
		resp, err := s.Client.Do(req)
		s.Assert().NoError(err)
		s.Assert().Equal(http.StatusBadRequest, resp.StatusCode)
		respByte, err := ioutil.ReadAll(resp.Body)
		s.Assert().NoError(err)
		defer resp.Body.Close()

		respMap := map[string]interface{}{}
		err = json.Unmarshal(respByte, &respMap)
		s.Assert().NoError(err)
		s.Assert().NotNil(respMap)
		s.Assert().Equal("API_VALIDATION_ERROR", respMap["error_code"])
	})

	s.T().Run("Return 400 for missing required field", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"lastName": "Doe",
			"address":  "Singapore",
			"isActive": true,
		}
		jbyt, err := json.Marshal(reqBody)
		s.Assert().NoError(err)

		req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(jbyt))
		s.Assert().NoError(err)

		req.Header = RequiredHeaders
		resp, err := s.Client.Do(req)
		s.Assert().NoError(err)
		s.Assert().Equal(http.StatusBadRequest, resp.StatusCode)
		respByte, err := ioutil.ReadAll(resp.Body)
		s.Assert().NoError(err)
		defer resp.Body.Close()

		respMap := map[string]interface{}{}
		err = json.Unmarshal(respByte, &respMap)
		s.Assert().NoError(err)
		s.Assert().NotNil(respMap)
		s.Assert().Equal("API_VALIDATION_ERROR", respMap["error_code"])
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
	s.Assert().NoError(err)

	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(jbyt))
	s.Assert().NoError(err)

	req.Header = RequiredHeaders
	resp, err := s.Client.Do(req)
	s.Assert().NoError(err)
	s.Assert().Equal(http.StatusCreated, resp.StatusCode)
	respByte, err := ioutil.ReadAll(resp.Body)
	s.Assert().NoError(err)
	defer resp.Body.Close()

	respMap := map[string]interface{}{}
	err = json.Unmarshal(respByte, &respMap)
	s.Assert().NoError(err)
	s.assertUserValue(reqBody, respMap)
	// new request for GET
	getUserByIDURI := fmt.Sprintf("%s/api/users/%v", s.Host, respMap["id"])
	req, err = http.NewRequest(http.MethodGet, getUserByIDURI, nil)
	s.Assert().NoError(err)

	req.Header = RequiredHeaders
	resp, err = s.Client.Do(req)
	s.Assert().NoError(err)
	s.Assert().Equal(http.StatusOK, resp.StatusCode)
	respByte, err = ioutil.ReadAll(resp.Body)
	s.Assert().NoError(err)
	defer resp.Body.Close()

	respMap = map[string]interface{}{}
	err = json.Unmarshal(respByte, &respMap)
	s.Assert().NoError(err)
	s.assertUserValue(reqBody, respMap)

}

func (s userTestSuite) TestUpdateUserUsingPUT() {
	uri := fmt.Sprintf("%s/api/users", s.Host)
	reqBody := map[string]interface{}{
		"firstName": "John",
		"lastName":  "Doe",
		"address":   "Singapore",
		"isActive":  true,
	}
	jbyt, err := json.Marshal(reqBody)
	s.Assert().NoError(err)
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(jbyt))
	s.Assert().NoError(err)

	req.Header = RequiredHeaders
	resp, err := s.Client.Do(req)
	s.Assert().NoError(err)
	s.Assert().Equal(http.StatusCreated, resp.StatusCode)
	respByte, err := ioutil.ReadAll(resp.Body)
	s.Assert().NoError(err)
	defer resp.Body.Close()

	respMap := map[string]interface{}{}
	err = json.Unmarshal(respByte, &respMap)
	s.Assert().NoError(err)
	s.assertUserValue(reqBody, respMap)

	// new request for GET
	getUserByIDURI := fmt.Sprintf("%s/api/users/%v", s.Host, respMap["id"])
	req, err = http.NewRequest(http.MethodGet, getUserByIDURI, nil)
	s.Assert().NoError(err)

	req.Header = RequiredHeaders
	resp, err = s.Client.Do(req)
	s.Assert().NoError(err)
	s.Assert().Equal(http.StatusOK, resp.StatusCode)
	respByte, err = ioutil.ReadAll(resp.Body)
	s.Assert().NoError(err)
	defer resp.Body.Close()

	respMap = map[string]interface{}{}
	err = json.Unmarshal(respByte, &respMap)
	s.Assert().NoError(err)
	s.assertUserValue(reqBody, respMap)

	updateUserByIDURI := fmt.Sprintf("%s/api/users/%v", s.Host, respMap["id"])
	reqBodyUpdate := map[string]interface{}{
		"firstName": "John",
		"lastName":  "Doe",
		"address":   "Jakarta",
		"isActive":  true,
	}
	jbyt, err = json.Marshal(reqBodyUpdate)
	s.Assert().NoError(err)

	req, err = http.NewRequest(http.MethodPut, updateUserByIDURI, bytes.NewBuffer(jbyt))
	s.Assert().NoError(err)

	req.Header = RequiredHeaders
	resp, err = s.Client.Do(req)
	s.Assert().NoError(err)
	s.Assert().Equal(http.StatusOK, resp.StatusCode)
	respByte, err = ioutil.ReadAll(resp.Body)
	s.Assert().NoError(err)
	defer resp.Body.Close()

	updateRespMap := map[string]interface{}{}
	err = json.Unmarshal(respByte, &updateRespMap)
	s.Assert().NoError(err)
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
	s.Assert().NoError(err)

	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(jbyt))
	s.Assert().NoError(err)

	req.Header = RequiredHeaders
	resp, err := s.Client.Do(req)
	s.Assert().NoError(err)
	s.Assert().Equal(http.StatusCreated, resp.StatusCode)
	respByte, err := ioutil.ReadAll(resp.Body)
	s.Assert().NoError(err)
	defer resp.Body.Close()

	respMap := map[string]interface{}{}
	err = json.Unmarshal(respByte, &respMap)
	s.Assert().NoError(err)
	s.assertUserValue(reqBody, respMap)
	// new request for GET: assert the inserted User is exists
	getUserByIDURI := fmt.Sprintf("%s/api/users/%v", s.Host, respMap["id"])
	req, err = http.NewRequest(http.MethodGet, getUserByIDURI, nil)
	s.Assert().NoError(err)

	req.Header = RequiredHeaders
	resp, err = s.Client.Do(req)
	s.Assert().NoError(err)
	s.Assert().Equal(http.StatusOK, resp.StatusCode)
	respByte, err = ioutil.ReadAll(resp.Body)
	s.Assert().NoError(err)
	defer resp.Body.Close()

	respMap = map[string]interface{}{}
	err = json.Unmarshal(respByte, &respMap)
	s.Assert().NoError(err)
	s.assertUserValue(reqBody, respMap)

	// new request for DELETE: assert the inserted User is exists
	deleteUserByIDURI := fmt.Sprintf("%s/api/users/%v", s.Host, respMap["id"])
	req, err = http.NewRequest(http.MethodDelete, deleteUserByIDURI, nil)
	s.Assert().NoError(err)

	req.Header = RequiredHeaders
	resp, err = s.Client.Do(req)
	s.Assert().NoError(err)
	s.Assert().Equal(http.StatusNoContent, resp.StatusCode)
	closerDelete := ioutil.NopCloser(resp.Body)
	s.Assert().NoError(err)
	defer closerDelete.Close()

	// new request for GET: assert the delete action is success User is exists
	getDeletedUserByIDURI := fmt.Sprintf("%s/api/users/%v", s.Host, respMap["id"])
	req, err = http.NewRequest(http.MethodGet, getDeletedUserByIDURI, nil)
	s.Assert().NoError(err)

	req.Header = RequiredHeaders
	resp, err = s.Client.Do(req)
	s.Assert().NoError(err)
	s.Assert().Equal(http.StatusNotFound, resp.StatusCode)
	closerGet := ioutil.NopCloser(resp.Body)
	s.Assert().NoError(err)
	defer closerGet.Close()
}

func (s *userTestSuite) assertUserValue(expected map[string]interface{}, actual map[string]interface{}) {
	s.Assert().NotNil(actual)
	s.Assert().NotNil(actual["id"])
	s.Assert().NotEmpty(actual["id"])
	s.Assert().Equal(expected["firstName"], actual["firstName"])
	s.Assert().Equal(expected["lastName"], actual["lastName"])
	s.Assert().Equal(expected["address"], actual["address"])
	s.Assert().Equal(expected["isActive"], actual["isActive"])
}
