package youcrawl

import (
	"fmt"
	"io"
	"net/http"
)

// make request with url
func RequestWithURL(task *Task, middlewares ...Middleware) (io.Reader, error) {
	req, err := http.NewRequest("GET", task.Url, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	if task.Context.Cookie != nil {
		client.Jar = task.Context.Cookie
	}
	for _, middleware := range middlewares {
		middleware.Process(client, req, &task.Context)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	EngineLogger.Info(fmt.Sprintf("%s [%d]", task.Url, resp.StatusCode))
	task.Context.Request = req
	task.Context.Response = resp

	for _, middleware := range middlewares {
		middleware.RequestCallback(client, req, &task.Context)
	}
	return resp.Body, nil
}
