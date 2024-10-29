package command

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrInvalidSyntax = errors.New("invalid syntax")
	ErrInvalidID     = errors.New("id is invalid or could not be parsed")
	ErrFailedToBuild = errors.New("failed to build command")
)

type Keyword string

var (
	Len     Keyword = "len"
	Push    Keyword = "push"
	Pop     Keyword = "pop"
	Drain   Keyword = "drain"
	Consume Keyword = "consume"
	Stop    Keyword = "stop"
)

type Command struct {
	ID      uuid.UUID
	Keyword Keyword
	Args    []string
}

func Build(keyword Keyword, args ...string) (Command, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return Command{}, ErrFailedToBuild
	}

	return Command{
		ID:      id,
		Keyword: keyword,
		Args:    args,
	}, nil
}

func Parse(input string) (Command, error) {
	spl := strings.Split(input, "::")

	if len(spl) < 2 {
		return Command{}, ErrInvalidSyntax
	}

	id, err := uuid.Parse(spl[0])
	if err != nil {
		return Command{}, fmt.Errorf("%w: %w", ErrInvalidID, err)
	}

	cmd := Command{
		ID:      id,
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
	case Consume:
		if len(cmd.Args) > 0 {
			return cmd, fmt.Errorf("%w: consume takes no args", ErrInvalidSyntax)
		}
	case Stop:
		if len(cmd.Args) > 0 {
			return cmd, fmt.Errorf("%w: stop takes no args", ErrInvalidSyntax)
		}
	default:
		return Command{ID: id}, fmt.Errorf("%w: unknown keyword", ErrInvalidSyntax)
	}

	return cmd, nil
}

func (c Command) String() string {
	out := fmt.Sprintf("%s::%s", c.ID.String(), string(c.Keyword))
	for _, a := range c.Args {
		out = fmt.Sprintf("%s::%s", out, a)
	}
	return out
}
