package youcrawl

import (
	"sync"
	"testing"
)

func TestProxy_UseMiddleware(t *testing.T) {
	e := NewEngine(&EngineOption{MaxRequest: 3})
	e.AddURLs("https://www.example.com")
	e.AddHTMLParser(DefaultTestParser)
	e.UseMiddleware(ProxyMiddleware)
	var wg sync.WaitGroup
	wg.Add(1)
	e.Run(&wg)
	wg.Wait()
}
