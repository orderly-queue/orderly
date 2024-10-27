package command

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidSyntax = errors.New("invalid syntax")
)

type Keyword string

var (
	Len   Keyword = "len"
	Push  Keyword = "push"
	Pop   Keyword = "pop"
	Drain Keyword = "drain"
)

type Command struct {
	ID      string
	Keyword Keyword
	Args    []string
}

func Parse(input string) (Command, error) {
	spl := strings.Split(input, "::")

	if len(spl) < 2 {
		return Command{ID: spl[0]}, ErrInvalidSyntax
	}

	cmd := Command{
		ID:      spl[0],
		Keyword: Keyword(spl[1]),
		Args:    []string{},
	}

	if len(spl) > 2 {
		cmd.Args = spl[2:]
	}

	switch Keyword(spl[1]) {
	case Len:
		if len(cmd.Args) > 0 {
			return cmd, fmt.Errorf("%w: len takes no args", ErrInvalidSyntax)
		}
	case Push:
		if len(cmd.Args) == 0 {
			return cmd, fmt.Errorf("%w: push requires an arg", ErrInvalidSyntax)
		}
	case Pop:
		if len(cmd.Args) > 0 {
			return cmd, fmt.Errorf("%w: pop takes no args", ErrInvalidSyntax)
		}
	case Drain:
		if len(cmd.Args) > 0 {
			return cmd, fmt.Errorf("%w: drain takes no args", ErrInvalidSyntax)
		}
	default:
		return Command{ID: spl[0]}, fmt.Errorf("%w: unknown keyword", ErrInvalidSyntax)
	}

	return cmd, nil
}
