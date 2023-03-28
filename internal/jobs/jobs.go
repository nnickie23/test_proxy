package jobs

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/nnickie23/test_proxy/internal/entities/constants"
	"github.com/nnickie23/test_proxy/internal/entities/errors"
	"github.com/nnickie23/test_proxy/internal/entities/models"
	"github.com/nnickie23/test_proxy/internal/logger"
	"github.com/nnickie23/test_proxy/internal/service"
)

type St struct {
	logger  logger.Logger
	service service.Service

	stop   bool
	stopMU sync.Mutex
	wg     sync.WaitGroup
}

var taskCh = make(chan models.Task)
var resultCh = make(chan models.Task, 5)

func New(logger logger.Logger, service service.Service) *St {
	return &St{
		logger:  logger,
		service: service,
	}
}

func (j *St) Start(done context.Context) {
	j.wg.Add(1)
	// run goroutine which checks for existing tasks
	go j.checkTasks(done)
	j.wg.Add(1)
	// run goroutine which saves results of task in DB
	go j.saveTaskResults(done)
	// create workers
	j.createWorkers(10, done)
}

func (j *St) beginJob() bool {
	j.stopMU.Lock()
	defer j.stopMU.Unlock()
	return !(j.stop)
}

func (j *St) endJob() {
	j.wg.Done()
}

func (j *St) Stop() {
	j.stopMU.Lock()
	j.stop = true
	j.stopMU.Unlock()
	j.wg.Wait()
}

func (j *St) checkTasks(done context.Context) {
	defer j.endJob()
	var err error
	var tasks []*models.Task
	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-ticker.C:
			if !j.beginJob() {
				return
			}
			tasks, err = j.service.TaskGetNewSetInProcess(context.Background())
			if err != nil && err != errs.ObjectNotFound {
				j.logger.Errorw("Failed to list tasks", err)
			}
			if len(tasks) > 0 {
				for _, task := range tasks {
					taskCh <- *task
				}
			}
		case <-done.Done():
			ticker.Stop()
			return
		}
	}
}

func (j *St) createWorkers(n int, done context.Context) {
	for i := 0; i < n; i++ {
		j.wg.Add(1)
		// goroutine which executes tasks - http request to 3rd party services
		go j.makeRequest(done)
	}
}

func (j *St) makeRequest(done context.Context) {
	defer j.endJob()
	for {
		select {
		case task := <-taskCh:
			if !j.beginJob() {
				return
			}
			result, err := j.Request(task)
			if err != nil {
				j.logger.Errorw("", err)
			}
			resultCh <- result
		case <-done.Done():
			return
		}
	}
}

func (j *St) saveTaskResults(done context.Context) {
	defer j.endJob()
	counter := 0
	maxSleepTime := 60 * time.Second
	for {
		select {
		case r := <-resultCh:
			if !j.beginJob() {
				return
			}
			// perform operations with result entity r
			if err := j.service.TaskAddResults(context.Background(), &r); err != nil {
				j.logger.Errorw("", err)
				return
			}
			counter = 0
		case <-done.Done():
			return
		default:
			counter++
			if counter > 50 {
				sleepTime := time.Duration(counter) * 10 * time.Second
				if sleepTime > maxSleepTime {
					sleepTime = maxSleepTime
				}
				time.Sleep(sleepTime)
			} else {
				time.Sleep(10 * time.Second)
			}
		}
	}
}

func (j *St) Request(task models.Task) (models.Task, error) {
	result := models.Task{}

	request, err := http.NewRequest(*task.RequestMethod, *task.RequestUrl, nil)
	if err != nil {
		j.logger.Errorw("", err)
		return models.Task{}, err
	}
	request.Header.Add("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 60 * time.Second,
	}
	response, err := client.Do(request)
	if err != nil {
		j.logger.Errorw("", err)
		result.Status = &constants.Error
		return result, err
	}

	result.Uuid = task.Uuid
	result.Status = &constants.Done
	result.ResponseHttpStatusCode = &response.StatusCode

	var responseHeaders json.RawMessage
	if len(response.Header) == 0 {
		responseHeaders = json.RawMessage("{}")
	} else {
		responseHeadersMap := make(map[string]string)
		for k, v := range response.Header {
			responseHeadersMap[k] = v[0]
		}

		responseHeadersJson, err := json.Marshal(responseHeadersMap)
		if err != nil {
			j.logger.Errorw("", err)
			result.Status = &constants.Error
			return result, err
		}

		responseHeaders = json.RawMessage(responseHeadersJson)
	}
	result.ResponseHeaders = &responseHeaders

	var responseContentLength int64
	if response.Header.Get("Content-Length") == "" {
		responseContentLength = 0
	} else {
		responseContentLength, err = strconv.ParseInt(response.Header.Get("Content-Length"), 10, 64)
		if err != nil {
			j.logger.Errorw("", err)
			result.Status = &constants.Error
			return result, err
		}
	}
	result.ResponseContentLength = &responseContentLength

	return result, nil
}
