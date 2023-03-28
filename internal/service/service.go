package service

import (
	"github.com/nnickie23/test_proxy/internal/logger"
	"github.com/nnickie23/test_proxy/internal/repository"
)

type Service interface {
	Task
}

type service struct {
	logger     logger.Logger
	repository repository.Repository
}

func New(logger logger.Logger, repository repository.Repository) *service {
	return &service{
		logger: logger,
		repository: repository,
	}
}
