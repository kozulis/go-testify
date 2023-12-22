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

func TestMainHandlerWhenRequestIsNotValid(t *testing.T) {
	cases := []struct {
		requestString string
		statusCode    int
		err           string
		name          string
	}{
		{
			requestString: "/cafe?count=5&city=thula",
			statusCode:    http.StatusBadRequest,
			err:           "wrong city value",
			name:          "wrong city value",
		},
		{
			requestString: "/cafe?count=5",
			statusCode:    http.StatusBadRequest,
			err:           "wrong city value",
			name:          "city is empty",
		},
		{
			requestString: "/cafe?count=abc5&city=moscow",
			statusCode:    http.StatusBadRequest,
			err:           "wrong count value",
			name:          "wrong count value",
		},
		{
			requestString: "/cafe?city=moscow",
			statusCode:    http.StatusBadRequest,
			err:           "count missing",
			name:          "count is empty",
		},
		{
			requestString: "/cafe?count=0city=moscow",
			statusCode:    http.StatusBadRequest,
			err:           "wrong count value",
			name:          "count is 0",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			body, resp, err := testRequestConfig("GET", tc.requestString)

			require.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			assert.Equal(t, tc.err, string(body))
		})
	}
}

func TestMainHandlerWhenRequestIsValid(t *testing.T) {
	cases := []struct {
		count         int
		requestString string
		statusCode    int
		name          string
	}{
		{
			count:         1,
			requestString: "/cafe?count=1&city=moscow",
			statusCode:    http.StatusOK,
			name:          "1 city statusOk",
		},
		{
			count:         3,
			requestString: "/cafe?count=3&city=moscow",
			statusCode:    http.StatusOK,
			name:          "3 cities statusOk",
		},
		{
			count:         4,
			requestString: "/cafe?count=5&city=moscow",
			statusCode:    http.StatusOK,
			name:          "5 cities statusOk",
		},
		{
			count:         4,
			requestString: "/cafe?count=55&city=moscow",
			statusCode:    http.StatusOK,
			name:          "55 cities statusOk",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			body, resp, err := testRequestConfig("GET", tc.requestString)
			res := strings.Split(string(body), ",")

			require.NoError(t, err)
			assert.NotEmpty(t, body)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			assert.Len(t, res, tc.count)
		})
	}
}

func testRequestConfig(requestMethod string, path string) ([]byte, *http.Response, error) {
	req := httptest.NewRequest(requestMethod, path, nil)
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)
	resp := responseRecorder.Result()
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return body, resp, err
}
