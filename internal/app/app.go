package app

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/orderly-queue/orderly/internal/config"
	"github.com/orderly-queue/orderly/internal/crypto"
	"github.com/orderly-queue/orderly/internal/jwt"
	"github.com/orderly-queue/orderly/internal/metrics"
	"github.com/orderly-queue/orderly/internal/probes"
	"github.com/orderly-queue/orderly/internal/queue"
	"github.com/orderly-queue/orderly/internal/storage"
	"github.com/thanos-io/objstore"
)

type server interface {
	Start(context.Context) error
	Stop(context.Context) error
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Routes() []*echo.Route
}

type App struct {
	Version string

	Config *config.Config

	Http server

	Jwt *jwt.Jwt

	Queue *queue.Queue

	Probes  *probes.Probes
	Metrics *metrics.Metrics

	Encryption *crypto.Encrptor

	Storage objstore.Bucket
}

func New(ctx context.Context, conf *config.Config) (*App, error) {
	enc, err := crypto.NewEncryptor(conf.EncryptionKey)
	if err != nil {
		return nil, err
	}

	app := &App{
		Config: conf,

		Jwt: jwt.New(conf.JwtSecret),

		Queue: queue.New(),

		Encryption: enc,

		Probes:  probes.New(conf.Probes.Port),
		Metrics: metrics.New(conf.Telemetry.Metrics.Port),
	}

	if conf.Storage.Enabled {
		storage, err := storage.New(conf.Storage)
		if err != nil {
			return nil, err
		}
		app.Storage = storage
	}

	return app, nil
}
