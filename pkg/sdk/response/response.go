package response

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrInvalidFormat = errors.New("the response is in an unrecognisable format")
	ErrInvalidID     = errors.New("id could not be parsed or is invalid")
)

type Response struct {
	ID      uuid.UUID
	Message string
	Error   error
}

func (r Response) Err() error {
	return r.Error
}

func Build(id uuid.UUID, msg string) Response {
	return Response{
		ID:      id,
		Message: msg,
	}
}

func BuildError(id uuid.UUID, err error) Response {
	return Response{
		ID:    id,
		Error: err,
	}
}

func Error(id uuid.UUID, err error) Response {
	return BuildError(id, err)
}

func (r Response) String() string {
	if r.Error != nil {
		return fmt.Sprintf("%s::error%s", r.ID, r.Error.Error())
	}
	return fmt.Sprintf("%s::%s", r.ID, r.Message)
}

func Parse(resp string) (Response, error) {
	spl := strings.Split(resp, "::")
	if len(spl) != 2 && len(spl) != 3 {
		fmt.Println(resp)
		return Response{}, ErrInvalidFormat
	}

	id, err := uuid.Parse(spl[0])
	if err != nil {
		return Response{}, fmt.Errorf("%w: %w", ErrInvalidID, err)
	}

	if spl[1] == "error" {
		return BuildError(id, fmt.Errorf("%s", spl[2])), nil
	}
	return Build(id, spl[1]), nil
}
