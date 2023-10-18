package config

import (
	"context"
	"io"
	"time"

	"gopkg.in/yaml.v3"
)

type key int

const (
	ctxKey key = iota
)

type Config struct {
	Logger struct {
		Level string
	}
	Server struct {
		Host, Port   string
		ReadTimeout  time.Duration `yaml:"readTimeout"`
		WriteTimeout time.Duration `yaml:"writeTimeout"`
	}
	DB struct {
		Dsn           string
		MigrationsDir string `yaml:"migrationsDir"`
	}
	Storage string
}

func NewConfig(configFile io.Reader) (*Config, error) {
	config := &Config{}

	yamlDecoder := yaml.NewDecoder(configFile)
	if err := yamlDecoder.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKey, c)
}

func GetFromContext(ctx context.Context) *Config {
	return ctx.Value(ctxKey).(*Config)
}
