package repository

import (
	"context"

	"github.com/nnickie23/test_proxy/internal/entities/constants"
	"github.com/nnickie23/test_proxy/internal/entities/models"
	"github.com/nnickie23/test_proxy/internal/postgres"
)

type Task interface {
	TaskCreate(context context.Context, task *models.Task) error
	TaskGetStatus(context context.Context, uuid string) (*models.Task, error)
	TaskGetNewSetInProcess(context context.Context) ([]*models.Task, error)
	TaskAddResults(context context.Context, task *models.Task) error
}

func (d *repository) TaskCreate(context context.Context, task *models.Task) error {
	var err error

	context, err = d.database.ContextWithTransaction(context)
	if err != nil {
		return err
	}
	defer func() { d.database.RollbackContextTransaction(context) }()

	var result int64
	fields := d.taskGetFields(task)
	if err := d.database.InsertReturnId(&result, context, constants.Tasks, fields); err != nil {
		return err
	}

	err = d.database.CommitContextTransaction(context)
	if err != nil {
		return err
	}

	d.logger.Infof("created task with id: %d", result)

	return nil
}

func (d *repository) TaskGetStatus(context context.Context, uuid string) (*models.Task, error) {
	var err error

	context, err = d.database.ContextWithTransaction(context)
	if err != nil {
		return nil, err
	}
	defer func() { d.database.RollbackContextTransaction(context) }()

	result := models.Task{}
	fields := postgresdb.Fields{"uuid": uuid}
	err = d.database.Get(&result, context, constants.Tasks, fields)
	if err != nil {
		return nil, err
	}

	err = d.database.CommitContextTransaction(context)
	if err != nil {
		return nil, err
	}

	return &models.Task{
		Uuid:                   result.Uuid,
		Status:                 result.Status,
		ResponseHttpStatusCode: result.ResponseHttpStatusCode,
		ResponseHeaders:        result.ResponseHeaders,
		ResponseContentLength:  result.ResponseContentLength,
	}, nil
}

func (d *repository) TaskGetNewSetInProcess(context context.Context) ([]*models.Task, error) {
	var err error

	context, err = d.database.ContextWithTransaction(context)
	if err != nil {
		return nil, err
	}
	defer func() { d.database.RollbackContextTransaction(context) }()

	result := []models.Task{}
	field := "status"
	err = d.database.UpdateFieldFromTo(&result, context, constants.Tasks, field, constants.New, constants.InProcess)
	if err != nil {
		return nil, err
	}

	err = d.database.CommitContextTransaction(context)
	if err != nil {
		return nil, err
	}

	ret := make([]*models.Task, 0)
	for _, i := range result {
		ret = append(ret, &models.Task{
			Uuid:           i.Uuid,
			RequestMethod:  i.RequestMethod,
			RequestUrl:     i.RequestUrl,
			RequestHeaders: i.RequestHeaders,
		})
	}

	return ret, nil
}

func (d *repository) TaskAddResults(context context.Context, task *models.Task) error {
	var err error

	context, err = d.database.ContextWithTransaction(context)
	if err != nil {
		return err
	}
	defer func() { d.database.RollbackContextTransaction(context) }()

	fields := d.taskGetFields(task)
	uuid := *task.Uuid
	if err := d.database.UpdateByUuid(context, constants.Tasks, uuid, fields); err != nil {
		return err
	}

	err = d.database.CommitContextTransaction(context)
	if err != nil {
		return err
	}

	d.logger.Infof("updated task with uuid: %s", uuid)

	return nil
}

func (d *repository) taskGetFields(obj *models.Task) map[string]interface{} {
	ret := make(map[string]interface{})

	if obj.Id != nil {
		ret["id"] = *obj.Id
	}

	if obj.Uuid != nil {
		ret["uuid"] = *obj.Uuid
	}

	if obj.Status != nil {
		ret["status"] = *obj.Status
	}

	if obj.RequestMethod != nil {
		ret["request_method"] = *obj.RequestMethod
	}

	if obj.RequestUrl != nil {
		ret["request_url"] = *obj.RequestUrl
	}

	if obj.RequestHeaders != nil {
		ret["request_headers"] = *obj.RequestHeaders
	}

	if obj.ResponseHttpStatusCode != nil {
		ret["response_http_status_code"] = *obj.ResponseHttpStatusCode
	}

	if obj.ResponseHeaders != nil {
		ret["response_headers"] = *obj.ResponseHeaders
	}

	if obj.ResponseContentLength != nil {
		ret["response_content_length"] = *obj.ResponseContentLength
	}

	return ret
}
