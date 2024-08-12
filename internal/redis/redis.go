package redis

import (
	"github.com/henrywhitaker3/go-template/internal/config"
	"github.com/henrywhitaker3/go-template/internal/tracing"
	"github.com/redis/rueidis"
	"github.com/redis/rueidis/rueidisotel"
)

func New(conf *config.Config) (rueidis.Client, error) {
	var client rueidis.Client
	var err error
	opts := rueidis.ClientOption{
		InitAddress:   []string{conf.Redis.Addr},
		Password:      conf.Redis.Password,
		MaxFlushDelay: conf.Redis.MaxFlushDelay,
	}
	if conf.Telemetry.Tracing.Enabled {
		client, err = rueidisotel.NewClient(
			opts, rueidisotel.WithTracerProvider(tracing.TracerProvider),
		)
	} else {
		client, err = rueidis.NewClient(opts)
	}

	return client, err
}
