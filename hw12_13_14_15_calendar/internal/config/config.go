package config

import (
	"io"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Logger struct {
		Level string
	}
	Server struct {
		Host, Port string
	}
	Db struct {
		Dsn           string
		MigrationsDir string `yaml:"migrations-dir"`
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
