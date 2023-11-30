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
	Storage string
	App     struct {
		Scheduler struct {
			Period time.Duration
			Queue  string
		}
	}
	HTTPServer struct {
		Host, Port   string
		ReadTimeout  time.Duration `yaml:"readTimeout"`
		WriteTimeout time.Duration `yaml:"writeTimeout"`
	} `yaml:"httpServer"`
	GRPCServer struct {
		Host, Port     string
		ConnectTimeout time.Duration `yaml:"connectTimeout"`
	} `yaml:"grpcServer"`
	DB struct {
		Dsn           string
		MigrationsDir string `yaml:"migrationsDir"`
	}
	RabbitMQServer struct {
		Host, Port, Login, Password string
	} `yaml:"rabbitMQServer"`
}

func New(configFile io.Reader) (*Config, error) {
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
