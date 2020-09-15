package youcrawl

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"os"
	"sync"
	"testing"
)

type PrintGlobalStorePostProcess struct {
}

func (p *PrintGlobalStorePostProcess) Process(store GlobalStore) error {

	rawItems := store.GetValue("items")
	if rawItems != nil {
		items := rawItems.([]interface{})
		fmt.Println(fmt.Sprintf("total crawl %d items", len(items)))

	}
	return nil
}

func TestPostProcess(t *testing.T) {
	e := NewEngine(&EngineOption{MaxRequest: 2})
	urls := []string{"https://example.com", "https://example.com", "https://example.com"}
	e.AddURLs(urls...)
	e.AddHTMLParser(func(doc *goquery.Document, ctx *Context) error {
		item := ctx.Item.(DefaultItem)
		title := doc.Find("title").Text()
		item.SetValue("title", title)
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

func TestOutputCSVPostProcess_Process(t *testing.T) {
	store := &MemoryGlobalStore{
		Content: map[string]interface{}{},
	}
	items := make([]map[string]interface{}, 0)
	for idx := 0; idx < 10; idx++ {
		item := make(map[string]interface{})
		item["title"] = fmt.Sprintf("title %d", idx)
		item["content"] = fmt.Sprintf("content %d", idx)
		item["ignore"] = fmt.Sprintf("ignore %d", idx)
		if idx%2 == 0 {
			item["exist"] = true
		}
		items = append(items, item)
	}
	store.SetValue("items", items)
	postprocess := NewOutputCSVPostProcess(OutputCSVPostProcessOption{
		OutputPath: "./output.csv",
		WithHeader: true,
		Keys:       []string{"title", "content", "exist"},
		KeysMapping: map[string]string{
			"title": "webTitle",
			"exist": "webExist",
		},
		NotExistValue: "Undefined",
	})

	defer os.Remove("./output.csv")
	err := postprocess.Process(store)
	if err != nil {
		t.Error(err)
	}

}
