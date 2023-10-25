package log

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
)

type Handler struct {
	logger      logger.Logger
	nextHandler http.Handler
}

func NewHandler(logger logger.Logger, next http.Handler) http.Handler {
	return &Handler{
		logger:      logger,
		nextHandler: next,
	}
}

func (l *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	lrw := NewLoggingResponseWriter(writer)

	start := time.Now()
	l.nextHandler.ServeHTTP(lrw, request)
	end := time.Since(start)

	logJSON, err := json.Marshal(
		struct {
			IP        string
			Datetime  string
			Method    string
			Path      string
			HTTP      string
			Status    string
			Time      string
			UserAgent string
		}{
			IP:        request.RemoteAddr,
			Datetime:  time.Now().Format(time.RFC822),
			Method:    request.Method,
			Path:      request.URL.Path,
			HTTP:      request.Proto,
			Status:    strconv.Itoa(lrw.StatusCode),
			Time:      end.String(),
			UserAgent: request.UserAgent(),
		},
	)
	if err != nil {
		l.logger.Error(err.Error())
	}

	l.logger.Info(string(logJSON))
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
