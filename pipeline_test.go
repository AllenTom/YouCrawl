package youcrawl

import (
	"fmt"
	"os"
	"sync"
	"testing"
)

func TestEngine(t *testing.T) {
	e := NewEngine(&EngineOption{MaxRequest: 2})
	e.AddURLs("https://www.example.com")
	e.AddHTMLParser(func(ctx *Context) error {
		item := ctx.Item.(DefaultItem)
		doc := ctx.Doc
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
		GetUrls: func(item interface{}, store GlobalStore) []string {
			return []string{"https://golang.google.cn/lib/godoc/images/home-gopher.png"}
		},
		GetSaveFileName: func(item interface{}, store GlobalStore, rawURL string) string {
			return "downloadImage.png"
		},
		OnImageDownloadComplete: func(item interface{}, store GlobalStore, url string, downloadFilePath string) {
			fmt.Println(url)
			fmt.Println(downloadFilePath)
		},
		OnDone: func(item interface{}, store GlobalStore) {
			fmt.Println("all image downloaded")
		},
	}
	err := downloadPipeline.Process(
		DefaultItem{},
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
