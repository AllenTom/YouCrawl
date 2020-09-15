package youcrawl

import (
	"github.com/PuerkitoBio/goquery"
	"os"
	"sync"
	"testing"
)

func TestEngine(t *testing.T) {
	e := NewEngine(&EngineOption{MaxRequest: 2})
	e.AddURLs("https://www.example.com")
	e.AddHTMLParser(func(doc *goquery.Document, ctx *Context) error {
		item := ctx.Item.(DefaultItem)
		title := doc.Find("title").Text()
		item.SetValue("title", title)

		return nil
	})
	itemLogPipeline := &ItemLogPipeline{
		Options: ItemLogPipelineOption{
			PrintTitle:   true,
			PrintDivider: true,
		},
	}
	e.AddPipelines(itemLogPipeline)
	e.UseMiddleware(&UserAgentMiddleware{})
	var wg sync.WaitGroup
	wg.Add(1)
	e.Run(&wg)
	wg.Wait()
}

func TestImageDownloadPipeline_Process(t *testing.T) {
	gb := MemoryGlobalStore{}
	downloadPipeline := ImageDownloadPipeline{
		GetStoreFileFolder: func(item interface{}, store GlobalStore) string {
			return "./download/crawl"
		},
		MaxDownload: 2,
		Middlewares: []Middleware{
			&UserAgentMiddleware{},
		},
	}
	err := downloadPipeline.Process(
		ImageDownloadItem{
			Urls: []string{"https://github.com/AllenTom/YouCrawl/raw/master/other/workflow.png"},
		},
		&gb,
	)
	if err != nil {
		t.Error(err)
	}
	defer func() {
		err := os.RemoveAll("./download")
		if err != nil {
			t.Error(err)
		}
	}()
}
