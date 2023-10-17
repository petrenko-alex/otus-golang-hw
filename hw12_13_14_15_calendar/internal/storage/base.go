package storage

import (
	"errors"
	"time"

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
	Update(entity.Event) error
	Delete(string) error
	GetAll() (*entity.Events, error)
	GetById(string) (*entity.Event, error)
	GetForPeriod(time.Time, time.Time) (*entity.Events, error)
	GetForTime(time.Time) (*entity.Event, error)
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
