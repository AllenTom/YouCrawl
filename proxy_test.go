package youcrawl

import (
	"sync"
	"testing"
)

func TestProxy_UseMiddleware(t *testing.T) {
	proxy, err := NewProxyMiddleware(ProxyMiddlewareOption{
		ProxyList: []string{"0.0.0.0"},
	})
	if err != nil {
		t.Error(err)
	}
	e := NewEngine(&EngineOption{MaxRequest: 3})
	e.AddURLs("https://www.example.com", "https://www.example.com", "https://www.example.com")
	e.AddHTMLParser(DefaultTestParser)
	e.UseMiddleware(proxy)
	var wg sync.WaitGroup
	wg.Add(1)
	e.Run(&wg)
	wg.Wait()
}
