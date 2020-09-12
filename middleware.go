package youcrawl

import (
	"net/http"
)

type Middleware interface {
	Process(c *http.Client, r *http.Request, ctx *Context)
	RequestCallback(c *http.Client, r *http.Request, ctx *Context)
}
