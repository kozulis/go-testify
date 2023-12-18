package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4
	body, _, err := testRequestSettings("GET", "/cafe?count=55&city=moscow")
	res := strings.Split(string(body), ",")

	require.NoError(t, err)
	assert.NotEmpty(t, body)
	assert.Len(t, res, totalCount)
}

func TestMainHandlerWhenCityIsNotAllowed(t *testing.T) {
	body, resp, err := testRequestSettings("GET", "/cafe?count=2&city=tver")
	wrongCityErrorString := "wrong city value"

	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, wrongCityErrorString, string(body))
}

func TestMainHandlerWhenRequestIsValid(t *testing.T) {
	body, resp, err := testRequestSettings("GET", "/cafe?count=5&city=moscow")

	require.NoError(t, err)
	assert.NotEmpty(t, body)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func testRequestSettings(requestMethod string, path string) ([]byte, *http.Response, error) {
	req := httptest.NewRequest(requestMethod, path, nil)
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)
	resp := responseRecorder.Result()
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return body, resp, err
}
