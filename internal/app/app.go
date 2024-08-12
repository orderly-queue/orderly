package app

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/henrywhitaker3/go-template/database/queries"
	"github.com/henrywhitaker3/go-template/internal/config"
	"github.com/henrywhitaker3/go-template/internal/jwt"
	"github.com/henrywhitaker3/go-template/internal/metrics"
	"github.com/henrywhitaker3/go-template/internal/postgres"
	"github.com/henrywhitaker3/go-template/internal/probes"
	"github.com/henrywhitaker3/go-template/internal/redis"
	"github.com/henrywhitaker3/go-template/internal/users"
	"github.com/redis/rueidis"
)

type server interface {
	Start(context.Context) error
	Stop(context.Context) error
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type App struct {
	Version string

	Config *config.Config

	Http server

	Probes  *probes.Probes
	Metrics *metrics.Metrics

	Users *users.Users

	Jwt *jwt.Jwt

	Database *sql.DB
	Queries  *queries.Queries
	Redis    rueidis.Client
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

	app := &App{
		Config: conf,

		Database: db,
		Queries:  queries,
		Redis:    redis,

		Users: users.New(queries),

		Jwt: jwt.New(conf.JwtSecret, redis),

		Probes:  probes.New(conf.Probes.Port),
		Metrics: metrics.New(conf.Telemetry.Metrics.Port),
	}

	return app, nil
}
