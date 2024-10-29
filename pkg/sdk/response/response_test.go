package response_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/orderly-queue/orderly/pkg/sdk/response"
	"github.com/stretchr/testify/require"
)

func TestItParsesResponses(t *testing.T) {
	type testCase struct {
		name     string
		response string
		expected response.Response
		errors   bool
	}

	id := uuid.New()

	tcs := []testCase{
		{
			name:     "parses ok response",
			response: fmt.Sprintf("%s::ok", id.String()),
			expected: response.Response{
				ID:      id,
				Message: "ok",
			},
		},
		{
			name:     "parses nil response",
			response: fmt.Sprintf("%s::nil", id.String()),
			expected: response.Response{
				ID:      id,
				Message: "nil",
			},
		},
		{
			name:     "parses error response",
			response: fmt.Sprintf("%s::error::some error", id.String()),
			expected: response.Response{
				ID:    id,
				Error: fmt.Errorf("%s", "some error"),
			},
		},
		{
			name:     "errors when parsing invalid id format",
			response: "error::bongo",
			errors:   true,
		},
	}

	for _, c := range tcs {
		t.Run(c.name, func(t *testing.T) {
			resp, err := response.Parse(c.response)
			if c.errors {
				require.NotNil(t, err)
				return
			} else {
				require.Nil(t, err)
			}

			t.Log(resp)

			require.Equal(t, c.expected.ID, resp.ID)
			require.Equal(t, c.expected.Message, resp.Message)
			require.Equal(t, c.expected.Error, resp.Error)
		})
	}
}
