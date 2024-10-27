package app

import (
	"context"
	"database/sql"
	"net/http"

	gocache "github.com/henrywhitaker3/go-cache"
	"github.com/orderly-queue/orderly/database/queries"
	"github.com/orderly-queue/orderly/internal/config"
	"github.com/orderly-queue/orderly/internal/crypto"
	"github.com/orderly-queue/orderly/internal/jwt"
	"github.com/orderly-queue/orderly/internal/metrics"
	"github.com/orderly-queue/orderly/internal/postgres"
	"github.com/orderly-queue/orderly/internal/probes"
	"github.com/orderly-queue/orderly/internal/redis"
	"github.com/orderly-queue/orderly/internal/storage"
	"github.com/orderly-queue/orderly/internal/users"
	"github.com/orderly-queue/orderly/internal/workers"
	"github.com/labstack/echo/v4"
	"github.com/redis/rueidis"
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

	Probes  *probes.Probes
	Metrics *metrics.Metrics

	Runner *workers.Runner

	Users *users.Users

	Jwt        *jwt.Jwt
	Encryption *crypto.Encrptor

	Database *sql.DB
	Queries  *queries.Queries
	Redis    rueidis.Client
	Storage  objstore.Bucket
	Cache    *gocache.Cache
}

func New(ctx context.Context, conf *config.Config) (*App, error) {
	redis, err := redis.New(conf)
	if err != nil {
		return nil, err
	}

	db, err := postgres.Open(ctx, conf.Database.Url, conf.Telemetry.Tracing)
	if err != nil {
		return nil, err
	}
	queries := queries.New(db)

	enc, err := crypto.NewEncryptor(conf.EncryptionKey)
	if err != nil {
		return nil, err
	}

	app := &App{
		Config: conf,

		Database: db,
		Queries:  queries,
		Redis:    redis,
		Cache:    gocache.NewCache(gocache.NewRueidisStore(redis)),

		Users: users.New(queries),

		Encryption: enc,
		Jwt:        jwt.New(conf.JwtSecret, redis),

		Probes:  probes.New(conf.Probes.Port),
		Metrics: metrics.New(conf.Telemetry.Metrics.Port),

		Runner: workers.NewRunner(ctx, redis),
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
