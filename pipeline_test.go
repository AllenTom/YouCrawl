package youcrawl

import (
	"github.com/PuerkitoBio/goquery"
	"sync"
	"testing"
)

func TestEngine(t *testing.T) {
	e := NewEngine(&EngineOption{MaxRequest: 2})
	e.AddURLs("https://www.example.com")
	e.AddHTMLParser(func(doc *goquery.Document, ctx Context) error {
		title := doc.Find("title").Text()
		ctx.Item.SetValue("title", title)
		return nil
	})
	itemLogPipeline := &ItemLogPipeline{
		Options: ItemLogPipelineOption{
			PrintTitle:   true,
			PrintDivider: true,
		},
	}
	e.AddPipelines(itemLogPipeline)
	e.UseMiddleware(UserAgentMiddleware)
	var wg sync.WaitGroup
	wg.Add(1)
	e.Run(&wg)
	wg.Wait()
}
