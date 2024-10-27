package connect

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/olahol/melody"
	"github.com/orderly-queue/orderly/internal/app"
	"github.com/orderly-queue/orderly/internal/command"
	"github.com/orderly-queue/orderly/internal/queue"
)

type ConnectHandler struct {
	app *app.App
}

func NewConnect(app *app.App) *ConnectHandler {
	return &ConnectHandler{app: app}
}

func (h *ConnectHandler) Handler() echo.HandlerFunc {
	m := melody.New()
	m.Config.MaxMessageSize = 256000
	m.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	return func(c echo.Context) error {
		m.HandleMessage(func(s *melody.Session, b []byte) {
			cmd, err := command.Parse(string(b))
			if err != nil {
				fail(s, cmd.ID, err)
				return
			}

			switch cmd.Keyword {
			case command.Len:
				h.len(s, cmd)
			case command.Push:
				h.push(s, cmd)
			case command.Pop:
				h.pop(s, cmd)
			}
		})

		return m.HandleRequest(c.Response(), c.Request())
	}
}

func (c *ConnectHandler) len(s *melody.Session, cmd command.Command) error {
	len := c.app.Queue.Len()
	return respond(s, command.Build(cmd.ID, fmt.Sprintf("%d", len)))
}

func (c *ConnectHandler) push(s *melody.Session, cmd command.Command) error {
	c.app.Queue.Push(cmd.Args[0])
	return respond(s, command.Build(cmd.ID, "ok"))
}

func (c *ConnectHandler) pop(s *melody.Session, cmd command.Command) error {
	item, err := c.app.Queue.Pop()
	if err != nil {
		if errors.Is(err, queue.ErrEmptyQueue) {
			return respond(s, command.Build(cmd.ID, "nil"))
		}
		return fail(s, cmd.ID, err)
	}
	return respond(s, command.Build(cmd.ID, item))
}

func respond(s *melody.Session, resp command.Response) error {
	return s.Write([]byte(resp.String()))
}

func fail(s *melody.Session, id string, err error) error {
	return s.Write([]byte(command.Build(id, fmt.Sprintf("error::%s", err.Error())).String()))
}

func (c *ConnectHandler) Method() string {
	return http.MethodGet
}

func (c *ConnectHandler) Path() string {
	return "/connect"
}

func (c *ConnectHandler) Middleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}
