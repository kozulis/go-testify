package main

import (
	"fmt"
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
	cases := []struct {
		count int
		city  string
		name  string
	}{
		{
			count: 4,
			city:  "moscow",
			name:  "valid 4 cities",
		},
		{
			count: 5,
			city:  "moscow",
			name:  "valid 5 cities",
		},
		{
			count: 55,
			city:  "moscow",
			name:  "valid 55 cities",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			body, _, err := testRequestSettings("GET", fmt.Sprintf("/cafe?count=%d&city=moscow", tc.count))
			res := strings.Split(string(body), ",")

			require.NoError(t, err)
			assert.NotEmpty(t, body)
			assert.Len(t, res, totalCount)
		})
	}
}

func TestMainHandlerWhenCityIsNotAllowed(t *testing.T) {
	wrongCityErrorString := "wrong city value"
	cases := []struct {
		count int
		city  string
		name  string
	}{
		{
			count: 4,
			city:  "thula",
			name:  "not valid thula",
		},
		{
			count: 5,
			city:  "tver",
			name:  "not valid tver",
		},
		{
			count: 55,
			city:  "lipetzk",
			name:  "not valid lipetzk",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			body, resp, err := testRequestSettings("GET", fmt.Sprintf("/cafe?count=%d&city=%s", tc.count, tc.city))

			require.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			assert.Equal(t, wrongCityErrorString, string(body))
		})
	}
}

func TestMainHandlerWhenRequestIsValid(t *testing.T) {
	cases := []struct {
		count int
		city  string
		name  string
	}{
		{
			count: 4,
			city:  "moscow",
			name:  "valid ok 4 moscow",
		},
		{
			count: 2,
			city:  "moscow",
			name:  "valid ok 2 moscow",
		},
		{
			count: 1,
			city:  "moscow",
			name:  "valid ok 1 moscow",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			body, resp, err := testRequestSettings("GET", fmt.Sprintf("/cafe?count=%d&city=moscow", tc.count))

			require.NoError(t, err)
			assert.NotEmpty(t, body)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})
	}
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
