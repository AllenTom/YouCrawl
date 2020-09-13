package youcrawl

import (
	"net/http"
)

type Middleware interface {
	// before request call
	Process(c *http.Client, r *http.Request, ctx *Context)
	// after request call
	RequestCallback(c *http.Client, r *http.Request, ctx *Context)
}
