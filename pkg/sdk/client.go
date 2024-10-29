package sdk

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/orderly-queue/orderly/pkg/sdk/command"
	"github.com/orderly-queue/orderly/pkg/sdk/response"
	"github.com/prometheus/client_golang/prometheus"
)

type MetricsConfig struct {
	RequestSecondHistogram *prometheus.HistogramVec
	RequestErrorCounter    *prometheus.CounterVec
	TransimitErrorCounter  prometheus.Counter
}

type ClientConfig struct {
	Endpoint string

	// The duration to wait whilst attempting to send a message
	SendTimeout time.Duration

	Metrics MetricsConfig
}

func (c ClientConfig) validate() error {
	if c.Endpoint == "" {
		return fmt.Errorf("%w: endpoint cannot be empty", ErrInvalidConfig)
	}
	return nil
}

func (c ClientConfig) URL() (*url.URL, error) {
	url, err := url.Parse(c.Endpoint)
	if err != nil {
		return nil, err
	}
	scheme := "wss"
	if url.Scheme == "http" {
		scheme = "ws"
	}
	url.Scheme = scheme
	url.Path = "/connect"
	return url, nil
}

type Client struct {
	listenMutex *sync.RWMutex
	pipes       map[uuid.UUID]chan response.Response
	tx          *sync.Mutex
	rx          chan string

	isClosed bool

	cancel    context.CancelFunc
	closed    chan struct{}
	closeOnce *sync.Once

	metrics MetricsConfig

	conn *websocket.Conn

	writeTimeout time.Duration
}

func NewClient(ctx context.Context, config ClientConfig) (*Client, error) {
	if err := config.validate(); err != nil {
		return nil, err
	}
	if config.SendTimeout == 0 {
		config.SendTimeout = time.Second * 3
	}

	ctx, cancel := context.WithCancel(ctx)

	c := &Client{
		listenMutex: &sync.RWMutex{},
		pipes:       make(map[uuid.UUID]chan response.Response),
		tx:          &sync.Mutex{},
		rx:          make(chan string, 100),
		closed:      make(chan struct{}, 1),
		closeOnce:   &sync.Once{},
		cancel:      cancel,

		metrics:      config.Metrics,
		writeTimeout: config.SendTimeout,
	}
	c.initMetrics()

	// Validate the websocket url
	url, err := config.URL()
	if err != nil {
		return nil, err
	}

	ws, _, err := websocket.DefaultDialer.DialContext(ctx, url.String(), make(http.Header))
	if err != nil {
		return nil, err
	}
	ws.SetCloseHandler(func(code int, text string) error {
		return c.Close()
	})
	c.conn = ws

	go c.loop(ctx)
	go c.read(ctx)

	return c, nil
}

func (c *Client) Len(ctx context.Context) (uint, error) {
	cmd, err := command.Build(command.Len)
	if err != nil {
		return 0, err
	}

	out, err := c.send(ctx, cmd)
	if err != nil {
		return 0, err
	}

	if err := out.Err(); err != nil {
		return 0, err
	}

	len, err := strconv.Atoi(out.Message)
	if err != nil {
		return 0, err
	}

	return uint(len), nil
}

func (c *Client) Push(ctx context.Context, data string) error {
	cmd, err := command.Build(command.Push, data)
	if err != nil {
		return err
	}
	out, err := c.send(ctx, cmd)
	if err != nil {
		return err
	}
	if err := out.Err(); err != nil {
		return fmt.Errorf("%w: %s", ErrFailedToPush, err)
	}
	return nil
}

func (c *Client) Pop(ctx context.Context) (string, error) {
	cmd, err := command.Build(command.Pop)
	if err != nil {
		return "", err
	}

	out, err := c.send(ctx, cmd)
	if err != nil {
		return "", err
	}
	if err := out.Err(); err != nil {
		return "", fmt.Errorf("%w: %w", ErrFailedToPop, err)
	}

	if out.Message == "nil" {
		return "", ErrQueueEmpty
	}

	return out.Message, nil
}

func (c *Client) Consume(ctx context.Context) (<-chan string, error) {
	cmd, err := command.Build(command.Consume)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToConsume, err)
	}

	c.listenMutex.Lock()
	resp := make(chan response.Response, 100)
	c.pipes[cmd.ID] = resp
	c.listenMutex.Unlock()

	if err := c.transmit(cmd.String()); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToSend, err)
	}

	select {
	case <-ctx.Done():
		stop := command.Command{ID: cmd.ID, Keyword: command.Stop}
		c.transmit(stop.String())
		return nil, ctx.Err()
	case ok := <-resp:
		if err := ok.Err(); err != nil {
			return nil, fmt.Errorf("%w: %s", ErrFailedToConsume, err)
		}
	}

	out := make(chan string, 100)
	go func() {
		defer c.ignore(cmd.ID)
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-resp:
				out <- msg.Message
			}
		}
	}()

	return out, nil
}

func (c *Client) send(ctx context.Context, cmd command.Command) (*response.Response, error) {
	defer c.ignore(cmd.ID)
	resp := c.listen(cmd.ID)

	start := time.Now()

	if err := c.transmit(cmd.String()); err != nil {
		return nil, err
	}

	timeout, cancel := context.WithTimeout(ctx, c.writeTimeout)
	defer cancel()

	var out response.Response
	select {
	case <-timeout.Done():
		return nil, fmt.Errorf("%w: %w", ErrFailedToSend, ctx.Err())
	case out = <-resp:
		// We don't to do anything here
	}
	dur := time.Since(start)

	if c.metrics.RequestSecondHistogram != nil {
		c.metrics.RequestSecondHistogram.With(prometheus.Labels{"method": string(cmd.Keyword)}).Observe(dur.Seconds())
	}

	if out.Err() != nil {
		if c.metrics.RequestErrorCounter != nil {
			c.metrics.RequestErrorCounter.With(prometheus.Labels{"method": string(cmd.Keyword)}).Inc()
		}
	}

	return &out, nil
}

func (c *Client) transmit(msg string) error {
	c.tx.Lock()
	defer c.tx.Unlock()
	if err := c.conn.WriteMessage(
		websocket.TextMessage,
		[]byte(msg),
	); err != nil {
		if c.metrics.TransimitErrorCounter != nil {
			c.metrics.TransimitErrorCounter.Inc()
		}
		return fmt.Errorf("%w: %w", ErrFailedToSend, err)
	}
	return nil
}

func (c *Client) loop(ctx context.Context) {
	defer c.Close()
	for {
		select {
		case <-c.closed:
			return
		case <-ctx.Done():
			return
		case msg := <-c.rx:
			resp, err := response.Parse(msg)
			if err != nil {
				continue
			}

			c.listenMutex.RLock()
			pipe, ok := c.pipes[resp.ID]
			c.listenMutex.RUnlock()
			if !ok {
				continue
			}
			pipe <- resp
		}
	}
}

func (c *Client) read(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-c.closed:
			return
		default:
			_, msg, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseGoingAway) || errors.Is(err, syscall.ECONNRESET) {
					return
				}
				return
			}
			c.rx <- string(msg)
		}
	}
}

func (c *Client) listen(id uuid.UUID) <-chan response.Response {
	c.listenMutex.Lock()
	defer c.listenMutex.Unlock()
	l := make(chan response.Response, 1)
	c.pipes[id] = l
	return l
}

func (c *Client) ignore(id uuid.UUID) {
	c.listenMutex.Lock()
	defer c.listenMutex.Unlock()
	delete(c.pipes, id)
}

func (c *Client) Close() error {
	c.closeOnce.Do(func() {
		c.isClosed = true
		close(c.closed)
		c.cancel()
		c.conn.Close()
	})
	return nil
}

func (c *Client) Done() <-chan struct{} {
	return c.closed
}

func (c *Client) initMetrics() {
	if c.metrics.RequestErrorCounter != nil {
		c.metrics.RequestErrorCounter.With(prometheus.Labels{"method": "len"}).Add(0)
		c.metrics.RequestErrorCounter.With(prometheus.Labels{"method": "push"}).Add(0)
		c.metrics.RequestErrorCounter.With(prometheus.Labels{"method": "pop"}).Add(0)
	}
}
