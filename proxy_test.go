package youcrawl

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"testing"
)

func TestProxy_UseMiddleware(t *testing.T) {
	e := NewEngine(&EngineOption{MaxRequest: 3})
	e.AddURLs("https://www.google.com")
	e.AddHTMLParser(func(doc *goquery.Document, ctx Context) {
		title := doc.Find("title").Text()
		fmt.Println(fmt.Sprintf("%s [%d]", ctx.Request.URL.String(), ctx.Response.StatusCode))
		fmt.Println(title)
	})
	e.UseMiddleware(ProxyMiddleware)
	stopChannel := make(chan struct{})
	e.Run(stopChannel)
	<-stopChannel
}
