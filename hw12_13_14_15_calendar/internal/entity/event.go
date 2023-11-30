package entity

import (
	"errors"
	"time"
)

var ErrEventNotFound = errors.New("event not found")

type Events []Event

type Event struct {
	ID          string
	Title       string
	DateTime    time.Time
	Description string
	Duration    string
	RemindTime  time.Time

	UserID int
}

type EventMsg struct {
	ID       string
	Title    string
	DateTime time.Time
	UserID   int
}

func (e Event) ToMsg() EventMsg {
	return EventMsg{
		ID:       e.ID,
		Title:    e.Title,
		DateTime: e.DateTime,
		UserID:   e.UserID,
	}
}
