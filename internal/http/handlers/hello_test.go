package handlers_test

import (
	"net/http"
	"testing"

	"github.com/henrywhitaker3/go-template/internal/test"
	"github.com/stretchr/testify/require"
)

func TestItSaysHelloName(t *testing.T) {
	type testCase struct {
		name   string
		url    string
		output string
		code   int
	}

	tcs := []testCase{
		{
			name:   "henry",
			url:    "/henry",
			output: "Hello henry!",
			code:   http.StatusOK,
		},
		{
			name:   "john",
			url:    "/john",
			output: "Hello john!",
			code:   http.StatusOK,
		},
	}

	app, cancel := test.App(t)
	defer cancel()

	for _, c := range tcs {
		t.Run(c.name, func(t *testing.T) {
			rec := test.Get(app, c.url, "")
			require.Equal(t, c.code, rec.Code)
			body := rec.Body.String()
			require.Equal(t, c.output, body)
		})
	}
}
