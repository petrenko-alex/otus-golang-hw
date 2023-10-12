package entity

import (
	"errors"
	"time"
)

// todo: use time.Duration type

var (
	ErrEventNotFound = errors.New("event not found")
)

type Events []Event

type Event struct {
	ID          string
	Title       string
	DateTime    time.Time
	Description string
	Duration    string
	RemindTime  string

	UserId int
}
