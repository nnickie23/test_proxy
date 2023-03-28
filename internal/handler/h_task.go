package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/nnickie23/test_proxy/internal/entities/models"
)

// @Summary Create Task
// @Tags task
// @Description  Create task
// @ID task
// @Accept json
// @Produce json
// @Param request body entities.TaskCreateRequest true "create task"
// @Success 200 {object} string "task ID"
// @Failure 400,404 {object} ErrRespSt
// @Failure 500 {object} ErrRespSt
// @Failure default {object} ErrRespSt
// @Router /task [post]
func (h *handler) hTaskCreate(w http.ResponseWriter, r *http.Request) {
	reqObj := struct {
		Method  string          `json:"method"`
		Url     string          `json:"url"`
		Headers json.RawMessage `json:"headers"`
	}{}
	if !uParseRequestJSON(w, r, &reqObj) {
		return
	}

	uuid := uuid.New().String()

	task := &models.Task{
		Uuid:           &uuid,
		RequestMethod:  &reqObj.Method,
		RequestUrl:     &reqObj.Url,
		RequestHeaders: &reqObj.Headers,
	}

	err := h.service.TaskCreate(r.Context(), task)
	if uHandleServiceErr(err, w) {
		return
	}

	uRespondJSON(w, struct {
		Id string `json:"id"`
	}{
		Id: uuid,
	})
}

// @Summary Get task status
// @Tags task
// @Description  Get task status
// @ID task status
// @Produce json
// @Param taskID path string true "task ID"
// @Success 200 {object} entities.ResultEntity "task result"
// @Failure 400,404 {object} ErrRespSt
// @Failure 500 {object} ErrRespSt
// @Failure default {object} ErrRespSt
// @Router /task/{taskID} [get]
func (h *handler) hTaskGetStatus(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["taskID"]

	res, err := h.service.TaskGetStatus(r.Context(), uuid)
	if uHandleServiceErr(err, w) {
		return
	}

	respObj := struct {
		Id             string          `json:"id"`
		Status         string          `json:"status"`
		HttpStatusCode int             `json:"http_status_code"`
		Headers        json.RawMessage `json:"headers"`
		Length         int64           `json:"length"`
	}{}

	respObj.Id = uuid

	if res.Status != nil {
		respObj.Status = *res.Status
	}

	if res.ResponseHttpStatusCode != nil {
		respObj.HttpStatusCode = *res.ResponseHttpStatusCode
	}

	if res.ResponseHeaders != nil {
		respObj.Headers = *res.ResponseHeaders
	}

	if res.ResponseContentLength != nil {
		respObj.Length = *res.ResponseContentLength
	}

	uRespondJSON(w, respObj)
}
