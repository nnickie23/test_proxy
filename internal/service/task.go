package service

import (
	"context"

	"github.com/nnickie23/test_proxy/internal/entities/models"
)

type Task interface {
	TaskCreate(context context.Context, task *models.Task) error
	TaskGetStatus(context context.Context, uuid string) (*models.Task, error)
	TaskGetNewSetInProcess(context context.Context) ([]*models.Task, error)
	TaskAddResults(context context.Context, task *models.Task) error
}

func (s *service) TaskCreate(context context.Context, task *models.Task) error {
	if err := s.repository.TaskCreate(context, task); err != nil {
		return err
	}
	return nil
}

func (s *service) TaskGetStatus(context context.Context, uuid string) (*models.Task, error) {
	res, err := s.repository.TaskGetStatus(context, uuid)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *service) TaskGetNewSetInProcess(context context.Context) ([]*models.Task, error) {
	res, err := s.repository.TaskGetNewSetInProcess(context)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *service) TaskAddResults(context context.Context, task *models.Task) error {
	if err := s.repository.TaskAddResults(context, task); err != nil {
		return err
	}
	return nil
}
