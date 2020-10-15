package youcrawl

import (
	"net/http"
	"testing"
)

func TestDelayMiddleware_Process(t *testing.T) {
	req, err := http.NewRequest("GET", "http://www.example.com", nil)
	if err != nil {
		t.Error(err)
	}
	delayMiddleware := DelayMiddleware{
		Min:   1,
		Max:   2,
	}
	delayMiddleware.Process(nil,req,nil)

	delayMiddleware2 := DelayMiddleware{
		Fixed: 1,
	}
	delayMiddleware2.Process(nil,req,nil)
}