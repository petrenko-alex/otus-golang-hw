package internalhttp

import (
	"net/http"
)

type AppHandler struct {
	app    Application
	logger Logger
}

func NewAppHandler(app Application, logger Logger) *AppHandler {
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
