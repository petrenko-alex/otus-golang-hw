package internalhttp

import (
	"net/http"

	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
)

type AppHandler struct {
	app    Application
	logger logger.Logger
}

func NewAppHandler(app Application, logger logger.Logger) *AppHandler {
	return &AppHandler{
		app:    app,
		logger: logger,
	}
}

func (h AppHandler) ServeHTTP(writer http.ResponseWriter, _ *http.Request) {
	_, err := writer.Write([]byte("Hello World"))
	if err != nil {
		h.logger.Error(err.Error())
	}
}
