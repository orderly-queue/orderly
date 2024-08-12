package test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/henrywhitaker3/go-template/internal/app"
	"github.com/henrywhitaker3/go-template/internal/config"
	"github.com/henrywhitaker3/go-template/internal/http"
	"github.com/henrywhitaker3/go-template/internal/logger"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	root string
)

func init() {
	re := regexp.MustCompile(`^(.*go-template)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))
	root = string(rootPath)
}

func App(t *testing.T) (*app.App, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	logger.Wrap(ctx, zap.NewAtomicLevelAt(zapcore.DebugLevel))
	pgCont, err := postgres.Run(
		ctx,
		"postgres:16",
		testcontainers.WithLogger(testcontainers.TestLogger(t)),
		postgres.WithDatabase("orderly"),
		postgres.WithUsername("orderly"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	require.Nil(t, err)

	conf, err := config.Load(fmt.Sprintf("%s/api.example.yaml", root))
	require.Nil(t, err)
	conn, err := pgCont.ConnectionString(context.Background())
	require.Nil(t, err)
	conf.Database.Url = conn

	redisCont, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Image:        "ghcr.io/dragonflydb/dragonfly:latest",
				ExposedPorts: []string{"6379/tcp"},
				WaitingFor:   wait.ForListeningPort("6379/tcp"),
			},
			Started: true,
			Logger:  testcontainers.TestLogger(t),
		},
	)
	require.Nil(t, err)
	redisHost, err := redisCont.Host(ctx)
	require.Nil(t, err)
	redisPort, err := redisCont.MappedPort(ctx, nat.Port("6379"))
	require.Nil(t, err)
	conf.Redis.Addr = fmt.Sprintf("%s:%d", redisHost, redisPort.Int())

	conf.Environment = "testing"

	app, err := app.New(ctx, conf)
	require.Nil(t, err)

	app.Http = http.New(app)

	return app, func() {
		require.Nil(t, redisCont.Terminate(ctx))
		require.Nil(t, pgCont.Terminate(ctx))
		cancel()
	}
}
