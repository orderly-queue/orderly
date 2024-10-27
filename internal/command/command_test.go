package command

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestItParsesCommand(t *testing.T) {
	type testCase struct {
		name     string
		input    string
		expected Command
		errors   bool
	}

	tcs := []testCase{
		{
			name:  "parses len command",
			input: "bongo::len",
			expected: Command{
				ID:      "bongo",
				Keyword: Len,
				Args:    []string{},
			},
		},
		{
			name:  "parses push command",
			input: "bongo::push::apple",
			expected: Command{
				ID:      "bongo",
				Keyword: Push,
				Args:    []string{"apple"},
			},
		},
		{
			name:  "parses pop command",
			input: "bongo::pop",
			expected: Command{
				ID:      "bongo",
				Keyword: Pop,
				Args:    []string{},
			},
		},
		{
			name:  "parses drain command",
			input: "bongo::drain",
			expected: Command{
				ID:      "bongo",
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
			}

			require.Equal(t, c.expected.ID, cmd.ID)
			require.Equal(t, c.expected.Keyword, cmd.Keyword)
			require.Equal(t, c.expected.Args, cmd.Args)
		})
	}
}
