package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nnickie23/test_proxy/internal/configs"
	"github.com/nnickie23/test_proxy/internal/handler"
	"github.com/nnickie23/test_proxy/internal/jobs"
	"github.com/nnickie23/test_proxy/internal/logger"
	"github.com/nnickie23/test_proxy/internal/postgres"
	"github.com/nnickie23/test_proxy/internal/repository"
	"github.com/nnickie23/test_proxy/internal/server"
	"github.com/nnickie23/test_proxy/internal/service"
)

func Run() {
	var err error

	configs, err := configs.Load()
	if err != nil {
		log.Fatal(err)
	}

	if configs.App.Debug {
		mux := http.NewServeMux()
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
		go func() { log.Println(http.ListenAndServe(":8888", mux)) }()
	}

	app := struct {
		logger     logger.Logger
		database   postgresdb.PostgresDb
		repository repository.Repository
		service    service.Service
		handler    handler.Handler
		server     *server.St
		jobs       *jobs.St
	}{}

	app.logger, err = logger.RegisterLog(configs.App.LogLevel, configs.App.Mode)
	if err != nil {
		log.Fatal(err)
	}

	app.database, err = postgresdb.Open(app.logger, configs.Postgres)
	if err != nil {
		app.logger.Fatal(err)
	}
	defer app.database.Close()

	app.repository = repository.New(app.logger, app.database)

	app.service = service.New(app.logger, app.repository)

	app.handler = handler.New(app.logger, app.service)

	app.jobs = jobs.New(app.logger, app.service)

	app.server = server.New(app.logger, app.handler, configs.Http)

	app.logger.Infow(
		"Starting",
		"http_addr", configs.Http.Addr,
	)

	app.server.Start()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	doneCtx, done := context.WithCancel(context.Background())

	app.jobs.Start(doneCtx)

	var exitCode int

	select {
	case <-stop:
		fmt.Println("shutting down...")
	case <-app.server.Wait():
		exitCode = 1
	}

	shutdownCtx, shutdown := context.WithTimeout(context.Background(), 20 * time.Second)
	defer shutdown()

	err = app.server.Shutdown(shutdownCtx)
	if err != nil {
		app.logger.Errorw("Fail to shutdown http-api", err)
		exitCode = 1
	}

	app.logger.Infow("Waiting running jobs...")
	done()
	app.jobs.Stop()

	os.Exit(exitCode)
}
