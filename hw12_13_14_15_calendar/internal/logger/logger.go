package logger

import (
	"io"
	"log"
)

type Level int

const (
	Debug Level = iota
	Info
	Warning
	Error
)

type Logger interface {
	Debug(string)
	Info(string)
	Warning(string)
	Error(string)
}

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

func New(level Level, dst io.Writer) Logger {
	return SimpleLogger{
		Level:  level,
		logger: log.New(dst, "", 0),
	}
}
