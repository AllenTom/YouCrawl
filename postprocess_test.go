package youcrawl

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"sync"
	"testing"
)

type PrintGlobalStorePostProcess struct {
}

func (p *PrintGlobalStorePostProcess) Process(store GlobalStore) error {

	rawItems := store.GetValue("items")
	if rawItems != nil {
		items := rawItems.([]map[string]interface{})
		fmt.Println(fmt.Sprintf("total crawl %d items", len(items)))

	}
	return nil
}

func TestPostProcess(t *testing.T) {
	e := NewEngine(&EngineOption{MaxRequest: 2})
	urls := []string{"https://example.com", "https://example.com", "https://example.com"}
	e.AddURLs(urls...)
	e.AddHTMLParser(func(doc *goquery.Document, ctx *Context) error {
		title := doc.Find("title").Text()
		ctx.Item.SetValue("title", title)
		return nil
	})
	e.UseMiddleware(&UserAgentMiddleware{})
	e.AddPipelines(&GlobalStorePipeline{})
	e.AddPostProcess(&PrintGlobalStorePostProcess{})
	var wg sync.WaitGroup
	wg.Add(1)
	e.Run(&wg)
	wg.Wait()
}
