package sdk

import "errors"

var (
	ErrInvalidConfig   = errors.New("invalid config")
	ErrFailedToSend    = errors.New("failed to send command")
	ErrFailedToPush    = errors.New("failed to push")
	ErrFailedToPop     = errors.New("failed to pop")
	ErrQueueEmpty      = errors.New("could not pop from empty queue")
	ErrFailedToConsume = errors.New("failed to consume")
	ErrClosed          = errors.New("client is closed")
)
