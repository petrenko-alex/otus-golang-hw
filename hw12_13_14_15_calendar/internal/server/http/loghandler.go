package internalhttp

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type LogHandler struct {
	logger      Logger
	nextHandler http.Handler
}

func NewLogHandler(logger Logger, next http.Handler) http.Handler {
	return &LogHandler{
		logger:      logger,
		nextHandler: next,
	}
}

func (l *LogHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	lrw := NewLoggingResponseWriter(writer)

	start := time.Now()
	l.nextHandler.ServeHTTP(lrw, request)
	end := time.Since(start)

	logJson, err := json.Marshal(
		struct {
			Ip        string
			Datetime  string
			Method    string
			Path      string
			Http      string
			Status    string
			Time      string
			UserAgent string
		}{
			Ip:        request.RemoteAddr,
			Datetime:  time.Now().Format(time.RFC822),
			Method:    request.Method,
			Path:      request.URL.Path,
			Http:      request.Proto,
			Status:    strconv.Itoa(lrw.StatusCode),
			Time:      end.String(),
			UserAgent: request.UserAgent(),
		},
	)
	if err != nil {
		l.logger.Error(err.Error())
	}

	l.logger.Info(string(logJson))
}

type LoggingResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{w, http.StatusOK}
}

func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.StatusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
