package redis

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/orderly-queue/orderly/internal/config"
	"github.com/redis/rueidis"
	"github.com/redis/rueidis/rueidisotel"
)

var (
	ErrLocked = errors.New("key already locked")
)

func New(conf *config.Config) (rueidis.Client, error) {
	opts := rueidis.ClientOption{
		InitAddress:   []string{conf.Redis.Addr},
		Password:      conf.Redis.Password,
		MaxFlushDelay: conf.Redis.MaxFlushDelay,
	}

	var client rueidis.Client
	var err error
	if conf.Telemetry.Tracing.Enabled {
		client, err = rueidisotel.NewClient(opts, rueidisotel.WithDBStatement(func(cmdTokens []string) string {
			return strings.Join(cmdTokens, " ")
		}))
	} else {
		client, err = rueidis.NewClient(opts)
	}

	return client, err
}

func Lock(ctx context.Context, client rueidis.Client, key string, exp time.Duration) error {
	cmd := client.B().Set().Key(fmt.Sprintf("locks:%s", key)).Value("true").Nx().Px(exp).Build()
	res := client.Do(ctx, cmd)
	if err := res.Error(); err != nil {
		if errors.Is(err, rueidis.Nil) {
			return ErrLocked
		}
		return err
	}
	return nil
}

func Unlock(ctx context.Context, client rueidis.Client, key string) error {
	cmd := client.B().Del().Key(fmt.Sprintf("locks:%s", key)).Build()
	return client.Do(ctx, cmd).Error()
}
