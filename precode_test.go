package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4
	req := httptest.NewRequest("GET", "/cafe?count=55&city=moscow", nil) // здесь нужно создать запрос к сервису
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)
	resp := responseRecorder.Result()
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	res := strings.Split(string(body), ",")

	assert.NotEmpty(t, body)
	assert.Equal(t, totalCount, len(res))
}

func TestMainHandlerWhenCityIsNotAllowed(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=2&city=tver", nil)
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)
	resp := responseRecorder.Result()
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	statusCodeBadRequest := 400
	wrongCityErrorString := "wrong city value"

	assert.Equal(t, statusCodeBadRequest, resp.StatusCode)
	assert.Equal(t, wrongCityErrorString, string(body))
}

func TestMainHandlerWhenRequestIsValid(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=5&city=moscow", nil)
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)
	resp := responseRecorder.Result()
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	statusCodeOk := 200

	assert.NotEmpty(t, body)
	assert.Equal(t, statusCodeOk, resp.StatusCode)
}
