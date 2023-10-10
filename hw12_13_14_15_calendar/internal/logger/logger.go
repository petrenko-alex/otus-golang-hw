package logger

import (
	"io"
	"log"

	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/app"
)

type Level int

const (
	Debug Level = iota
	Info
	Warning
	Error
)

type SimpleLogger struct {
	Level Level

	logger *log.Logger
}

func (l SimpleLogger) Info(msg string) {
	if l.Level < Info {
		return
	}

	l.log("INFO: ", msg)
}

func (l SimpleLogger) Error(msg string) {
	if l.Level > Error {
		return
	}

	l.log("ERROR: ", msg)
}

func (l SimpleLogger) Debug(msg string) {
	if l.Level > Debug {
		return
	}

	l.log("DEBUG: ", msg)
}

func (l SimpleLogger) Warning(msg string) {
	if l.Level > Warning {
		return
	}

	l.log("WARNING: ", msg)
}

func (l SimpleLogger) log(prefix, msg string) {
	l.logger.SetPrefix(prefix)
	l.logger.Println(msg)
}

func New(level Level, dst io.Writer) app.Logger {
	return SimpleLogger{
		Level:  level,
		logger: log.New(dst, "", 0),
	}
}
