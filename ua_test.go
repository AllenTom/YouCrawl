package youcrawl

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"sync"
	"testing"
)

func TestReadUserAgentListFile(t *testing.T) {
	_, err := ReadListFile("./ua.txt")
	if err != nil {
		t.Error(err)
	}
}
func TestUserAgentPool_GetUserAgent(t *testing.T) {
	var wg sync.WaitGroup
	for idx := 0; idx < 10; idx++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			UserAgents.GetUserAgent()

		}(&wg)
	}
	wg.Wait()

}

func TestEngine_UseUAMiddleware(t *testing.T) {
	e := NewEngine(&EngineOption{MaxRequest: 3})
	e.AddURLs("https://www.example.com", "https://www.example.com", "https://www.example.com")
	e.AddHTMLParser(func(doc *goquery.Document, ctx Context) {
		title := doc.Find("title").Text()
		fmt.Println(fmt.Sprintf("%s [%d]", ctx.Request.URL.String(), ctx.Response.StatusCode))
		fmt.Println(title)
	})
	e.UseMiddleware(UserAgentMiddleware)
	stopChannel := make(chan struct{})
	e.Run(stopChannel)
	<-stopChannel
}
