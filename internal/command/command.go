package command

import (
	"errors"
	"strings"
)

var (
	ErrInvalidCommand = errors.New("invalid commanf")
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
		return Command{}, ErrInvalidCommand
	}

	var key Keyword
	switch Keyword(spl[1]) {
	case Len:
		key = Len
	case Push:
		key = Push
	case Pop:
		key = Pop
	case Drain:
		key = Drain
	default:
		return Command{}, ErrInvalidCommand
	}

	cmd := Command{
		ID:      spl[0],
		Keyword: key,
		Args:    []string{},
	}

	if len(spl) > 2 {
		cmd.Args = spl[2:]
	}

	return cmd, nil
}
