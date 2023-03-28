package handler

import (
	"net/http"

	"github.com/nnickie23/test_proxy/internal/logger"
	"github.com/nnickie23/test_proxy/internal/service"
)

// Handler - interface for handling web requests.
type Handler interface {
	InitRoutes() http.Handler
}

// handler - structure for implementation of handler interface
type handler struct {
	logger  logger.Logger
	service service.Service
}

// NewHandler - function for creating new Handler instance. Receives service and returns Handler interface
func New(logger logger.Logger, service service.Service) *handler {
	return &handler{
		logger: logger,
		service: service,
	}
}
