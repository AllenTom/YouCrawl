package youcrawl

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
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
	var LogMiddleware = func(c *http.Client, r *http.Request, ctx Context) {
		fmt.Println(fmt.Sprintf("request : %s", r.URL.String()))
	}
	e := NewEngine(&EngineOption{MaxRequest: 3})
	e.AddURLs("https://www.example.com")
	e.AddHTMLParser(DefaultTestParser)
	e.UseMiddleware(LogMiddleware)
	stopChannel := make(chan struct{})
	e.Run(stopChannel)
	<-stopChannel
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

func (i *ItemLogPipeline) Process(item *Item) error {
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
