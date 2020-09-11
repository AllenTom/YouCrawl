package youcrawl

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"sync"
	"sync/atomic"
	"testing"
)

var DefaultTestParser HTMLParser = func(doc *goquery.Document, ctx Context) error {
	title := doc.Find("title").Text()
	//fmt.Println(fmt.Sprintf("%s [%d]", ctx.Request.URL.String(), ctx.Response.StatusCode))
	fmt.Println(ctx.Request.Header.Get("User-Agent"))
	fmt.Println(title)
	return nil
}

func TestRequestWithURL(t *testing.T) {
	_, err := RequestWithURL(&Task{
		Url: "https://www.example.com",
	})
	if err != nil {
		t.Error(err)
	}
}

func TestEngine_Run(t *testing.T) {
	e := NewEngine(&EngineOption{MaxRequest: 3})
	e.AddURLs("https://www.example.com", "https://www.example.com", "https://www.example.com")
	//e.AddHTMLParser(func(doc *goquery.Document, ctx Context) error {
	//	ctx
	//	return nil
	//})
	var wg sync.WaitGroup
	wg.Add(1)
	e.Run(&wg)
	wg.Wait()
}

func TestParseHTML(t *testing.T) {
	bodyReader, err := RequestWithURL(&Task{
		Url: "https://www.example.com",
	})
	if err != nil {
		t.Error(err)
	}
	err = ParseHTML(bodyReader, func(doc *goquery.Document, ctx Context) error {
		title := doc.Find("title").Text()
		fmt.Println(title)
		return nil
	}, Context{content: map[string]interface{}{}})
	if err != nil {
		t.Error(err)
	}
}

type ItemLogPipelineOption struct {
	PrintTitle   bool
	PrintDivider bool
}

type ItemLogPipeline struct {
	Options ItemLogPipelineOption
}

func (i *ItemLogPipeline) Process(item *Item, _ *GlobalStore) error {
	if i.Options.PrintTitle {
		title, err := item.GetString("title")
		if err != nil {
			return err
		}
		fmt.Println("=====================   " + title + "   =====================")
	}
	if i.Options.PrintDivider {
		fmt.Println("==============================================================")
	}
	return nil
}

//http://www.eeeeeeeeeeeeexaaaaaaaaaaaaaaample.com/
func TestWebNotReach(t *testing.T) {
	e := NewEngine(&EngineOption{MaxRequest: 3})
	e.AddURLs("http://www.eeeeeeeeeeeeexaaaaaaaaaaaaaaample.com")
	var wg sync.WaitGroup
	wg.Add(1)
	e.Run(&wg)
	wg.Wait()
}

func TestAddTaskInRun(t *testing.T) {
	var hasAdd int64 = 0
	e := NewEngine(&EngineOption{MaxRequest: 5})
	e.AddURLs("http://www.example.com")
	e.AddHTMLParser(func(doc *goquery.Document, ctx Context) error {
		if hasAdd == 0 {
			ctx.Pool.AddURLs("http://www.example.com")
			atomic.AddInt64(&hasAdd, 1)
		}
		return nil
	})
	var wg sync.WaitGroup
	wg.Add(1)
	e.Run(&wg)
	wg.Wait()
}
