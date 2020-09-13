# Middleware
`Middleware` is a component used to modify the request. The component can modify the request or store information before and after the request. For example, `Cookie`, `UserAgent`, etc. are all components that modify the request.

## Custom Middleware

Custom Middleware only needs to implement the `youcrawl.Middleware` interface, use `e.UseMiddleware(&yourMiddleware)`

```go
type Middleware interface {
	// before request call
	Process(c *http.Client, r *http.Request, ctx *Context)
	// after request call
	RequestCallback(c *http.Client, r *http.Request, ctx *Context)
}
```

## UserAgent Middleware
UserAgent is one of the built-in middleware. It will read the file `./ua.txt` in the directory (each line represents a UserAgent), and randomly select a UserAgent and add it to the `Header` of the `Request`.

## Proxy Middleware
The Proxy middleware is use for reading `./proxy.txt` in the directory (each line represents a proxy address), and will randomly select a proxy when requesting.

## Cookie Middleware
Cookie Middleware can store requested cookies, supports multiple cookies, similar to HashMap.
```go
cookieMiddleware := NewCookieMiddleware(
    CookieMiddlewareOption{
	// The strategy for obtaining the index key
	// if not set, all requests will use the same cookie
    GetKey: nil,
})
e.UseMiddleware(cookieMiddleware)
```