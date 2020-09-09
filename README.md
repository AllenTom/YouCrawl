# YouCrawl

![](https://img.shields.io/badge/-Go-black?logo=go)

Simple web crawl

## Feature
HTML parser : [PuerkitoBio/goquery](https://github.com/PuerkitoBio/goquery)
## Example
```go
func main() {
    // new engine
    e := youcrawl.NewEngine(&youcrawl.EngineOption{
        // Up to three tasks can be executed concurrently
        MaxRequest: 3
    })
    // add urls to crawl
    e.AddURLs("https://www.zhihu.com")
    // add HTML parser
	e.AddHTMLParser(func(doc *goquery.Document, ctx youcrawl.Context) {
        // get document and your code
        // go `goquery` doc to get more detail
        title := doc.Find("title").Text()
        
        fmt.Println(fmt.Sprintf("%s [%d]",ctx.Request.URL.String(),ctx.Response.StatusCode))
        
		fmt.Println(title)
    })
    // make channel for wait
    stopChannel := make(chan struct{})
    // run crawler
    e.Run(stopChannel)
    // wait for it done
	<-stopChannel
}
```
