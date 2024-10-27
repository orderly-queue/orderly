package command

import "fmt"

type Response struct {
	ID      string
	Message string
}

func Build(id string, msg string) Response {
	return Response{
		ID:      id,
		Message: msg,
	}
}

func (r Response) String() string {
	return fmt.Sprintf("%s::%s", r.ID, r.Message)
}
