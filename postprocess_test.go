package youcrawl

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"testing"
	"time"
)

type GlobalStorePipeline struct {
}

func (g *GlobalStorePipeline) Process(item *Item, store *GlobalStore) error {
	item.SetValue("time", time.Now().Format("2006-01-02 15:04:05"))
	store.Lock()
	defer store.Unlock()
	if store.Content["items"] == nil {
		store.Content["items"] = make([]map[string]interface{}, 0)
	}
	rawItems := store.Content["items"]
	items := rawItems.([]map[string]interface{})
	store.Content["items"] = append(items, item.Store)
	return nil
}

type PrintGlobalStorePostProcess struct {
}

func (p *PrintGlobalStorePostProcess) Process(store *GlobalStore) error {
	items := (store.Content["items"]).([]map[string]interface{})
	fmt.Println(fmt.Sprintf("total crawl %d items", len(items)))
	return nil
}

func TestPostProcess(t *testing.T) {
	e := NewEngine(&EngineOption{MaxRequest: 2})
	urls := []string{"https://example.com", "https://example.com", "https://example.com"}
	e.AddURLs(urls...)
	e.AddHTMLParser(func(doc *goquery.Document, ctx Context) error {
		title := doc.Find("title").Text()
		ctx.Item.SetValue("title", title)
		return nil
	})
	e.UseMiddleware(UserAgentMiddleware)
	e.AddPipelines(&GlobalStorePipeline{})
	e.AddPostProcess(&PrintGlobalStorePostProcess{})
	stopChannel := make(chan struct{})
	e.Run(stopChannel)
	<-stopChannel
}
