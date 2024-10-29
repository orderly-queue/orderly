package command

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestItParsesCommand(t *testing.T) {
	type testCase struct {
		name     string
		input    string
		expected Command
		errors   bool
	}

	id := uuid.New()

	tcs := []testCase{
		{
			name:  "parses len command",
			input: fmt.Sprintf("%s::len", id.String()),
			expected: Command{
				ID:      id,
				Keyword: Len,
				Args:    []string{},
			},
		},
		{
			name:  "parses push command",
			input: fmt.Sprintf("%s::push::apple", id.String()),
			expected: Command{
				ID:      id,
				Keyword: Push,
				Args:    []string{"apple"},
			},
		},
		{
			name:  "parses pop command",
			input: fmt.Sprintf("%s::pop", id.String()),
			expected: Command{
				ID:      id,
				Keyword: Pop,
				Args:    []string{},
			},
		},
		{
			name:  "parses drain command",
			input: fmt.Sprintf("%s::drain", id.String()),
			expected: Command{
				ID:      id,
				Keyword: Drain,
				Args:    []string{},
			},
		},
		{
			name:   "errors with too few parts",
			input:  "bongo",
			errors: true,
		},
		{
			name:   "errors invalid keyword",
			input:  "bongo::bongo",
			errors: true,
		},
	}

	for _, c := range tcs {
		t.Run(c.name, func(t *testing.T) {
			cmd, err := Parse(c.input)
			if c.errors {
				require.NotNil(t, err)
				return
			} else {
				require.Nil(t, err)
			}

			require.Equal(t, c.expected.ID, cmd.ID)
			require.Equal(t, c.expected.Keyword, cmd.Keyword)
			require.Equal(t, c.expected.Args, cmd.Args)
		})
	}
}
