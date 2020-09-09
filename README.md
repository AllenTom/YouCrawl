# YouCrawl

![](https://img.shields.io/badge/lang-Go-green)
![](https://travis-ci.com/AllenTom/YouCrawl.svg?branch=master)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FAllenTom%2FYouCrawl.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2FAllenTom%2FYouCrawl?ref=badge_shield)
[![codecov](https://codecov.io/gh/AllenTom/YouCrawl/branch/master/graph/badge.svg)](https://codecov.io/gh/AllenTom/YouCrawl)
[![BCH compliance](https://bettercodehub.com/edge/badge/AllenTom/YouCrawl?branch=master)](https://bettercodehub.com/)

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
    // use your custom middleware
    var LogMiddleware = func(r *http.Request,ctx Context) {
		fmt.Println(fmt.Sprintf("request : %s",r.URL.String()))
	}
    e.UseMiddleware(LogMiddleware)
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

## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FAllenTom%2FYouCrawl.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2FAllenTom%2FYouCrawl?ref=badge_large)
