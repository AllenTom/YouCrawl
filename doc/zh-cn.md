# YouCrawl

![](https://img.shields.io/badge/lang-Go-green)
![](https://travis-ci.com/AllenTom/YouCrawl.svg?branch=master)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FAllenTom%2FYouCrawl.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2FAllenTom%2FYouCrawl?ref=badge_shield)
[![codecov](https://codecov.io/gh/AllenTom/YouCrawl/branch/master/graph/badge.svg)](https://codecov.io/gh/AllenTom/YouCrawl)
[![BCH compliance](https://bettercodehub.com/edge/badge/AllenTom/YouCrawl?branch=master)](https://bettercodehub.com/)

[简体中文](doc/zh-cn.md) | [English](../README.md)

使用Go语言实现的爬虫库
## 安装
```
go get github.com/allentom/youcrawl
```
## 功能
HTML parser : [PuerkitoBio/goquery](https://github.com/PuerkitoBio/goquery)
## Workflow 工作原理
![](../other/workflow.png)

黄色部分会并行的执行
## 示例
```go

func main() {
    // 创建一个engine
    e := youcrawl.NewEngine(&youcrawl.EngineOption{
        // 最大同时请求数
        MaxRequest: 3
    })
    // 使用自定义middleware
    var LogMiddleware = func(r *http.Request,ctx Context) {
		fmt.Println(fmt.Sprintf("request : %s",r.URL.String()))
	}
    e.UseMiddleware(LogMiddleware)
    // 添加爬虫的网址
    e.AddURLs("https://www.zhihu.com")
    // HTML解析器
	e.AddHTMLParser(func(doc *goquery.Document, ctx youcrawl.Context) {
        // 详细的使用方法请参见 goquery
        title := doc.Find("title").Text()
        
        fmt.Println(fmt.Sprintf("%s [%d]",ctx.Request.URL.String(),ctx.Response.StatusCode))
        
		fmt.Println(title)
    })
    // 创建一个waitgroup等待完成
    var wg sync.WaitGroup
	wg.Add(1)
	e.Run(&wg)
	wg.Wait()
}
```
## Middleware
YouCrawl提供了一些内置的中间件，中间件是在执行请求之前进行对请求参数进行编辑的组件。

### UserAgent
随机UserAgent的Middleware会随机从UA池中随机选取UA，将UA的列表保存在`ua.txt`中，每一行代表一个

```go
func main() {
    ...
    e.UseMiddleware(youcrawl.UserAgentMiddleware)
}
```

### Proxy
随机代理会从`./proxy.txt`中读取列表并随机选择一个代理使用

```go
func main() {
    ...
    e.UseMiddleware(youcrawl.ProxyMiddleware)
}
```

## Pipelines
Pipeline是处理Item的组件

### 自定义 Pipeline
只需要实现`youcrawl.Pipeline`接口即可
```go
type MyCustomPipeline struct {

}

func (g *MyCustomPipeline) Process(item *youcrawl.Item, store *youcrawl.GlobalStore) error {
	item.SetValue("time",time.Now().Format("2006-01-02 15:04:05"))
	return nil
}
```
调用 `e.AddPipelines()`

### 存储Pipeline
全局储存拥有与engine相同的生命周期，在使用的过程中，可以将数据储存在这里

这个Pipeline会保存item至全局存储中名为`items`的列表中
```go
e.AddPipelines(&GlobalStorePipeline{})
```

## Post-Process（后期处理）
当所有的操作完成后，会执行后期处理，执行一些例如将Global Store中的数据一次性写入数据库，记录爬取状态记录

### 自定义 post-process
实现`youcrawl.PostProcess`接口
```go
type PrintGlobalStorePostProcess struct {}

func (p *PrintGlobalStorePostProcess) Process(store *youcrawl.GlobalStore) error {
	items := (store.Content["items"]).([]map[string]interface{})
	fmt.Println(fmt.Sprintf("total crawl %d items",len(items)))
	return nil
}
```
调用 ` e.AddPostProcess()` to add your custom post-process

### OutputJsonPostProcess
将 `GlobalStore["items"]`  中的Item写入json文件
## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FAllenTom%2FYouCrawl.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2FAllenTom%2FYouCrawl?ref=badge_large)
