package repository

import (
	"github.com/nnickie23/test_proxy/internal/logger"
	"github.com/nnickie23/test_proxy/internal/postgres"
)

type Repository interface {
	Task
}

type repository struct {
	logger   logger.Logger
	database postgresdb.PostgresDb
}

func New(logger logger.Logger, database postgresdb.PostgresDb) *repository {
	return &repository{
		logger: logger,
		database: database,
	}
}
