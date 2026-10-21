package api

import (
	"github.com/vucongthanh92/shortbit/internal/api/http"
	"github.com/vucongthanh92/shortbit/internal/share/application"
)

type ApiContainer struct {
	HttpServer *http.Server
	ViewWorker *application.ViewWorker
}

func NewApiContainer(httpServer *http.Server, viewWorker *application.ViewWorker) *ApiContainer {
	return &ApiContainer{HttpServer: httpServer, ViewWorker: viewWorker}
}
