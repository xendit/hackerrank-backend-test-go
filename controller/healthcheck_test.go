package controller_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/xendit/hackerrank-backend-test-go/controller"
)

func TestLiveness(t *testing.T) {
	e := echo.New()
	controller.InitHealthCheckHandler(e)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthcheck/liveness", nil)
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Result().StatusCode)
	expectedMapBody := map[string]interface{}{
		"status": "OK",
	}
	expectedJson, err := json.Marshal(expectedMapBody)
	require.NoError(t, err)
	require.JSONEq(t, string(expectedJson), rec.Body.String())
}

func TestReadiness(t *testing.T) {
	e := echo.New()
	controller.InitHealthCheckHandler(e)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthcheck/readiness", nil)
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Result().StatusCode)
	expectedMapBody := map[string]interface{}{
		"status": "OK",
	}
	expectedJson, err := json.Marshal(expectedMapBody)
	require.NoError(t, err)
	require.JSONEq(t, string(expectedJson), rec.Body.String())
}
