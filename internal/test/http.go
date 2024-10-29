package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/orderly-queue/orderly/internal/app"
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
	var reader io.Reader
	if body != nil {
		by, err := json.Marshal(body)
		require.Nil(t, err)
		reader = bytes.NewReader(by)
	}

	req := httptest.NewRequest(http.MethodPost, url, reader)
	if apikey != "" {
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", apikey))
	}
	req.Header.Set("Content-type", echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	app.Http.ServeHTTP(rec, req)

	return rec
}

func Patch(t *testing.T, app *app.App, url string, body any, apikey string) *httptest.ResponseRecorder {
	var reader io.Reader
	if body != nil {
		by, err := json.Marshal(body)
		require.Nil(t, err)
		reader = bytes.NewReader(by)
	}

	req := httptest.NewRequest(http.MethodPatch, url, reader)
	if apikey != "" {
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", apikey))
	}
	req.Header.Set("Content-type", echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	app.Http.ServeHTTP(rec, req)

	return rec
}

func Delete(t *testing.T, app *app.App, url string, body any, apikey string) *httptest.ResponseRecorder {
	var reader io.Reader
	if body != nil {
		by, err := json.Marshal(body)
		require.Nil(t, err)
		reader = bytes.NewReader(by)
	}

	req := httptest.NewRequest(http.MethodDelete, url, reader)
	if apikey != "" {
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", apikey))
	}
	req.Header.Set("Content-type", echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	app.Http.ServeHTTP(rec, req)

	return rec
}

// Creates a new httptest server pointing at the app, returns the url
func Server(app *app.App) (string, context.CancelFunc) {
	srv := httptest.NewServer(app.Http)

	return srv.URL, func() {
		srv.Close()
	}
}
