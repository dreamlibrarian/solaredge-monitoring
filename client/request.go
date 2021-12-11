package client

import (
	"fmt"
	"net/url"
	"time"

	"github.com/dreamlibrarian/solaredge-monitoring/api"
)

const (
	timeUnitParam  = "timeUnit"
	startTimeParam = "startTime"
	endTimeParam   = "endTime"
)

type Request struct {
	urlObject url.URL
	query     url.Values
}

func (c *Client) CreateRequest(path string) Request {
	var req Request
	req.urlObject = c.baseURL
	req.urlObject.Path = path

	req.query = make(url.Values)

	return req
}

func (c *Client) CreateRequestf(pathTemplate string, elements ...interface{}) Request {
	return c.CreateRequest(fmt.Sprintf(pathTemplate, elements...))
}

func (r *Request) SetParam(key, value string) *Request {
	r.query.Set(key, value)
	return r
}

func (r *Request) SetTimeParam(key string, value time.Time) *Request {
	vStr := api.ToTimestamp(value)
	r.query.Set(key, vStr)
	return r
}

func (r *Request) SetDateParam(key string, value time.Time) *Request {
	vStr := api.ToDatestamp(value)
	r.query.Set(key, vStr)
	return r
}

func (r *Request) SetTimeParams(timeUnit string, startTime, endTime time.Time) *Request {
	return r.SetParam(timeUnitParam, timeUnit).
		SetTimeParam(startTimeParam, startTime).
		SetTimeParam(endTimeParam, endTime)
}

func (r *Request) URL() *url.URL {
	u := r.urlObject
	u.RawQuery = r.query.Encode()

	return &u
}

func (r *Request) String() string {
	return r.URL().String()
}
