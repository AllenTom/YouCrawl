package youcrawl

import (
	"sync"
	"testing"
	"time"
)

// only two task can be run
func TestRequestPool_Close(t *testing.T) {
	e := NewEngine(&EngineOption{MaxRequest: 2})

	urls := []string{"https://example.com", "https://example.com", "https://example.com"}
	e.AddURLs(urls...)
	e.AddHTMLParser(func(ctx *Context) error {
		item := ctx.Item.(DefaultItem)
		doc := ctx.Doc
		title := doc.Find("title").Text()
		item.SetValue("title", title)
		<-time.After(3 * time.Second)
		return nil
	})
	go func() {
		<-time.After(1 * time.Second)
		e.InterruptChan <- struct{}{}
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	e.Run(&wg)
	wg.Wait()
}
