package main

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		count int
		want  int
	}{
		{count: 0, want: 0},
		{count: 1, want: 1},
		{count: 2, want: 2},
		{count: 100, want: len(cafeList["moscow"])},
	}

	for _, v := range requests {
		url := "/cafe?city=moscow&count=" + strconv.Itoa(v.count)
		rr := httptest.NewRecorder()
		req := httptest.NewRequest("GET", url, nil)
		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
		body := strings.TrimSpace(rr.Body.String())
		var cafes []string
		if body == "" {
			cafes = []string{}
		} else {
			cafes = strings.Split(body, ",")
		}
		assert.Equalf(t, v.want, len(cafes),
			"для count=%d ожидали %d кафе, получили %d: %v",
			v.count, v.want, len(cafes), cafes)
	}

}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		search    string
		wantCount int
	}{
		{search: "фасоль", wantCount: 0},
		{search: "кофе", wantCount: 2},
		{search: "вилка", wantCount: 1},
	}

	for _, v := range requests {
		url := "/cafe?city=moscow&search=" + v.search
		rr := httptest.NewRecorder()
		req := httptest.NewRequest("GET", url, nil)
		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)

		body := strings.TrimSpace(rr.Body.String())
		var cafes []string
		if body == "" {
			cafes = []string{}
		} else {
			cafes = strings.Split(body, ",")
		}

		assert.Equal(t, v.wantCount, len(cafes),
			"для search=%q ожидали %d кафе, получили %d: %v",
			v.search, v.wantCount, len(cafes), cafes)
	}
}
