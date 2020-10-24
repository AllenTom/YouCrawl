# Middleware
middleware是负责处理请求的组件，组件可以在请求之前以及之后对请求进行修改或者储存信息。例如Cookie、UserAgent等都是对request请求进行编辑与修改。

## 自定义Middleware

自定义Middleware只需要实现youcrawl.Middleware接口即可，使用`e.UseMiddleware(&yourMiddleware)`

```go
type Middleware interface {
	// before request call
	Process(c *http.Client, r *http.Request, ctx *Context)
	// after request call
	RequestCallback(c *http.Client, r *http.Request, ctx *Context)
}
```

## UserAgent Middleware
UserAgent是内置的中间件之一，会读取目录下的`./ua.txt`文件（每一行代表一个UserAgent），在请求时会随机选择一个UserAgent并添加至Request的Header中

## Proxy Middleware
Proxy 中间件负责读取目录下的`./proxy.txt`(每一行代表一个代理地址)，在请求时会随机选择一个代理。

## Cookie Middleware
Cookie Middleware可以储存请求过程中的Cookie，支持多Cookie使用，类似HashMap的。
```go
cookieMiddleware := NewCookieMiddleware(
    CookieMiddlewareOption{
    // 获取索引Key的策略，如果未设置，所有请求将会使用同一个Cookie记录
    GetKey: nil,
})
e.UseMiddleware(cookieMiddleware)
```

## Delay Middleware
在请求时，有可能需要延迟请求的需求，使得请求更加自然。Middleware支持随机产生也支持固定值，单位为秒。
```go
delayMiddleware := DelayMiddleware{
		Min:   1,
		Max:   2,
		Fixed: 1
	}
```