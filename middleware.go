package youcrawl

import (
	"net/http"
)

type Middleware func(r *http.Request, ctx Context)
