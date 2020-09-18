package youcrawl

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"sync"
	"testing"
)

var DefaultTestParser HTMLParser = func(ctx *Context) error {
	doc := ctx.Doc
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
	doc, err := goquery.NewDocumentFromReader(bodyReader)
	if err != nil {
		EngineLogger.Error(err)
	}
	err = ParseHTML(func(ctx *Context) error {
		title := doc.Find("title").Text()
		fmt.Println(title)
		return nil
	}, &Context{})
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

func (i *ItemLogPipeline) Process(item interface{}, _ GlobalStore) error {
	if i.Options.PrintTitle {
		item := item.(DefaultItem)
		title, _ := item.GetValue("title")
		fmt.Println("=====================   " + title.(string) + "   =====================")
	}
	if i.Options.PrintDivider {
		fmt.Println("==============================================================")
	}
	return nil
}

func TestCookie(t *testing.T) {
	e := NewEngine(&EngineOption{MaxRequest: 5})
	e.AddURLs("http://www.bing.com", "http://www.yandex.com")
	e.AddHTMLParser(DefaultTestParser)
	cookieMiddleware := NewCookieMiddleware(CookieMiddlewareOption{
		GetKey: nil,
	})
	e.UseMiddleware(cookieMiddleware)
	var wg sync.WaitGroup
	wg.Add(1)
	e.Run(&wg)
	wg.Wait()
}

func TestNewTask(t *testing.T) {
	e := NewEngine(&EngineOption{MaxRequest: 5})
	addTask := NewTask("http://www.bing.com", map[string]interface{}{
		"webLabel": "bing",
	})
	e.AddTasks(&addTask)
	e.AddHTMLParser(DefaultTestParser)
	var wg sync.WaitGroup
	wg.Add(1)
	e.Run(&wg)
	wg.Wait()
}
