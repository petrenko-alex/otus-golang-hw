package config

import (
	"errors"
	"io"

	"gopkg.in/yaml.v3"
)

var (
	ErrInvalidStorageValue = errors.New("invalid storage value in config")
)

type StorageType string

const (
	Memory StorageType = "memory"
	DB     StorageType = "db"
)

type Config struct {
	Logger struct {
		Level string
	}
	Server struct {
		Host, Port string
	}
	Db struct {
		Dsn string
	}
	Storage StorageType
}

func NewConfig(configFile io.Reader) (*Config, error) {
	config := &Config{}

	yamlDecoder := yaml.NewDecoder(configFile)
	if err := yamlDecoder.Decode(&config); err != nil {
		return nil, err
	}

	if validateErr := validateConfig(config); validateErr != nil {
		return nil, validateErr
	}

	return config, nil
}

func validateConfig(config *Config) error {
	if config.Storage != Memory && config.Storage != DB {
		return ErrInvalidStorageValue
	}

	return nil
}
