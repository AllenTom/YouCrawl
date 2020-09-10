package youcrawl

import (
	"testing"
)

func TestProxy_UseMiddleware(t *testing.T) {
	e := NewEngine(&EngineOption{MaxRequest: 3})
	e.AddURLs("https://www.example.com")
	e.AddHTMLParser(DefaultTestParser)
	e.UseMiddleware(ProxyMiddleware)
	stopChannel := make(chan struct{})
	e.Run(stopChannel)
	<-stopChannel
}
