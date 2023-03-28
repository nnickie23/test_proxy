package server

import (
	"context"
	"net/http"

	"github.com/nnickie23/test_proxy/internal/configs"
	"github.com/nnickie23/test_proxy/internal/handler"
	"github.com/nnickie23/test_proxy/internal/logger"
)

type St struct {
	logger logger.Logger
	lChan  chan error
	server *http.Server
}

func New(logger logger.Logger, handler handler.Handler, config configs.HttpConfig) *St {
	return &St{
		logger: logger,
		lChan: make(chan error, 1),
		server: &http.Server{
			Addr:         config.Addr,
			Handler:      handler.InitRoutes(),
			ReadTimeout:  config.ReadTimeout,
			WriteTimeout: config.WriteTimeout,
		},
	}
}

func (a *St) Start() {
	go func() {
		err := a.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			a.logger.Errorw("Http server closed", err)
			a.lChan <- err
		}
	}()
}

func (a *St) Wait() <-chan error {
	return a.lChan
}

func (a *St) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}
