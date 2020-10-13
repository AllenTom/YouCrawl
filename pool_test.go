package youcrawl

import (
	"fmt"
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
	e.RunAndWait()
}

func TestRequestPool_GetTotal(t *testing.T) {
	e := NewEngine(&EngineOption{MaxRequest: 2})
	urls := []string{"http://example.com", "https://example.com", "https://example.com"}
	e.AddURLs(urls...)
	e.AddHTMLParser(func(ctx *Context) error {
		item := ctx.Item.(DefaultItem)
		doc := ctx.Doc
		title := doc.Find("title").Text()
		item.SetValue("title", title)
		<-time.After(3 * time.Second)
		return nil
	})
	fmt.Println(e.Pool.GetTotal())
	fmt.Println(e.Pool.GetUnRequestCount())
	fmt.Println(e.Pool.GetCompleteCount())
}
