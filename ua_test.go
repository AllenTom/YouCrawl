package youcrawl

import (
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
	middleware, err := NewUserAgentMiddleware(UserAgentMiddlewareOption{
		UserAgentList: []string{
			"Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.130 Safari/537.36",
		},
	})
	if err != nil {
		t.Error(err)
	}
	var wg sync.WaitGroup
	for idx := 0; idx < 10; idx++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			middleware.GetUserAgent()
		}(&wg)
	}
	wg.Wait()

}

func TestEngine_UseUAMiddleware(t *testing.T) {
	middleware, err := NewUserAgentMiddleware(UserAgentMiddlewareOption{
		UserAgentList: []string{
			"Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.130 Safari/537.36",
		},
	})
	if err != nil {
		t.Error(err)
	}
	e := NewEngine(&EngineOption{MaxRequest: 3})
	e.AddURLs("https://www.example.com")
	e.AddHTMLParser(DefaultTestParser)
	e.UseMiddleware(middleware)
	var wg sync.WaitGroup
	wg.Add(1)
	e.Run(&wg)
	wg.Wait()
}
