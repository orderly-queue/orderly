package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/henrywhitaker3/go-template/internal/app"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func Get(app *app.App, url string, apikey string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, url, nil)
	if apikey != "" {
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", apikey))
	}
	rec := httptest.NewRecorder()
	app.Http.ServeHTTP(rec, req)

	return rec
}

func Post(t *testing.T, app *app.App, url string, body any, apikey string) *httptest.ResponseRecorder {
	by, err := json.Marshal(body)
	require.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(by))
	if apikey != "" {
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", apikey))
	}
	rec := httptest.NewRecorder()
	app.Http.ServeHTTP(rec, req)

	return rec
}

func Patch(t *testing.T, app *app.App, url string, body any, apikey string) *httptest.ResponseRecorder {
	by, err := json.Marshal(body)
	require.Nil(t, err)

	req := httptest.NewRequest(http.MethodPatch, url, bytes.NewReader(by))
	if apikey != "" {
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", apikey))
	}
	rec := httptest.NewRecorder()
	app.Http.ServeHTTP(rec, req)

	return rec
}

func Delete(t *testing.T, app *app.App, url string, body any, apikey string) *httptest.ResponseRecorder {
	by, err := json.Marshal(body)
	require.Nil(t, err)

	req := httptest.NewRequest(http.MethodDelete, url, bytes.NewReader(by))
	if apikey != "" {
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", apikey))
	}
	rec := httptest.NewRecorder()
	app.Http.ServeHTTP(rec, req)

	return rec
}
