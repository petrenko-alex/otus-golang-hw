package storage

import (
	"errors"

	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/entity"
	memorystorage "github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/storage/sql"
)

type Type string

const (
	Memory Type = "memory"
	DB     Type = "db"
)

var (
	ErrInvalidStorageValue = errors.New("invalid storage value in config")
)

type Storage interface {
	Create(entity.Event) (string, error)
	ReadOne(string) (*entity.Event, error)
	ReadAll() (*entity.Events, error)
	Update(entity.Event) error
	Delete(string) error
}

func GetStorage(storageType string) (Storage, error) {
	if Type(storageType) != Memory && Type(storageType) != DB {
		return nil, ErrInvalidStorageValue
	}

	if Type(storageType) == Memory {
		return memorystorage.New(), nil

	}

	return sqlstorage.New(), nil
}
