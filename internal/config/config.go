package config

import (
	"errors"
	"os"
	"time"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type LogLevel string

func (l LogLevel) Level() zap.AtomicLevel {
	switch l {
	case "debug":
		return zap.NewAtomicLevelAt(zap.DebugLevel)
	case "error":
		return zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "info":
		fallthrough
	default:
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	}
}

type Postgres struct {
	Url string `yaml:"url"`
}

type Redis struct {
	Addr          string        `yaml:"addr"`
	Password      string        `yaml:"password"`
	MaxFlushDelay time.Duration `yaml:"max_flush_delay"`
}

type Tracing struct {
	Enabled    bool    `yaml:"enabled"`
	SampleRate float64 `yaml:"sample_rate"`
	Endpoint   string  `yaml:"endpoint"`
}

type Metrics struct {
	Enabled bool `yaml:"enabled"`
	Port    int  `yaml:"port"`
}

type Telemetry struct {
	Tracing Tracing `yaml:"tracing"`
	Metrics Metrics `yaml:"metrics"`
}

type Probes struct {
	Port int `yaml:"port"`
}

type Http struct {
	Port int `yaml:"port"`
}

type Config struct {
	Environment string   `yaml:"env"`
	LogLevel    LogLevel `yaml:"log_level"`
	Database    Postgres `yaml:"database"`
	Redis       Redis    `yaml:"redis"`

	Probes Probes `yaml:"probes"`
	Http   Http   `yaml:"http"`

	Telemetry Telemetry `yaml:"telemetry"`
}

func Load(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var conf Config
	if err := yaml.Unmarshal(file, &conf); err != nil {
		return nil, err
	}

	conf.setDefaults()
	if err := conf.validate(); err != nil {
		return nil, err
	}

	return &conf, nil
}

func (c *Config) validate() error {
	if c.Database.Url == "" {
		return errors.New("invalid db url")
	}
	if c.Redis.Addr == "" {
		return errors.New("invalid redis addr")
	}
	return nil
}

func (c *Config) setDefaults() {
	if c.Environment == "" {
		c.Environment = "dev"
	}
	if c.Redis.MaxFlushDelay == 0 {
		c.Redis.MaxFlushDelay = time.Microsecond * 100
	}
	if c.Http.Port == 0 {
		c.Http.Port = 8765
	}
	if c.Telemetry.Metrics.Port == 0 {
		c.Telemetry.Metrics.Port = 8766
	}
	if c.Probes.Port == 0 {
		c.Probes.Port = 8767
	}
}
