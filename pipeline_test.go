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

func TestImageDownloadPipeline_Process(t *testing.T) {
	downloadPipeline := ImageDownloadPipeline{
		GetStoreFileFolder: func(item *Item, store *GlobalStore) string {
			return "./download/crawl"
		},
		MaxDownload: 2,
		Middlewares: []Middleware{
			UserAgentMiddleware,
		},
	}
	err := downloadPipeline.Process(
		&Item{
			Store: map[string]interface{}{
				"downloadImgURLs": []string{
					"https://www.flaticon.com/svg/static/icons/svg/3408/3408545.svg",
					"https://www.flaticon.com/svg/static/icons/svg/3408/3408540.svg",
					"https://www.flaticon.com/svg/static/icons/svg/3408/3408678.svg",
					"https://www.flaticon.com/svg/static/icons/svg/3408/3408736.svg",
				},
			},
		},
		&GlobalStore{},
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
