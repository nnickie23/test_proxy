package models

import "encoding/json"

type Task struct {
	Id                     *int64           `db:"id"`
	Uuid                   *string          `db:"uuid"`
	Status                 *string          `db:"status"`
	RequestMethod          *string          `db:"request_method"`
	RequestUrl             *string          `db:"request_url"`
	RequestHeaders         *json.RawMessage `db:"request_headers"`
	ResponseHttpStatusCode *int             `db:"response_http_status_code"`
	ResponseHeaders        *json.RawMessage `db:"response_headers"`
	ResponseContentLength  *int64           `db:"response_content_length"`
}
