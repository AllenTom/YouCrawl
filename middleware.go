package youcrawl

import (
	"net/http"
)

type Middleware func(c *http.Client, r *http.Request, ctx *Context)
