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
	for _, middleware := range middlewares {
		middleware(client, req, task.Context)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	fmt.Println(fmt.Sprintf("%s [%d]", task.Url, resp.StatusCode))
	task.Context.Request = req
	task.Context.Response = resp
	return resp.Body, nil
}
