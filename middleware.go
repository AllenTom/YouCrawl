package youcrawl

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
)

type Middleware interface {
	// before request call
	Process(c *http.Client, r *http.Request, ctx *Context)
	// after request call
	RequestCallback(c *http.Client, r *http.Request, ctx *Context)
}

type ProxyMiddleware struct {
	List []string
}

type ProxyMiddlewareOption struct {
	// set proxy list,
	// if both ProxyList and ProxyFilePath are provided,combine tow list
	ProxyList []string
	// read proxy from file,use `./proxy.txt` by default,
	// if both ProxyList and ProxyFilePath are provided,combine tow list
	ProxyFilePath string
}

func NewProxyMiddleware(option ProxyMiddlewareOption) (*ProxyMiddleware, error) {
	proxyList := option.ProxyList
	if proxyList == nil {
		proxyList = make([]string, 0)
	}

	list, err := readProxyListFile(option.ProxyFilePath)
	if err != nil {
		return nil, err
	}
	proxyList = append(proxyList, list...)
	middleware := &ProxyMiddleware{List: proxyList}
	return middleware, nil
}

func readProxyListFile(proxyPath string) ([]string, error) {
	if len(proxyPath) == 0 {
		proxyPath = "./proxy.txt"
	}
	rawList, err := ReadListFile(proxyPath)
	proxyList := make([]string, 0)
	// drop not valid url
	for _, rawUrl := range rawList {
		_, err = url.Parse(rawUrl)
		if err != nil {
			continue
		}
		proxyList = append(proxyList, rawUrl)
	}
	if err != nil {
		return nil, err
	}
	return proxyList, nil
}

func (p *ProxyMiddleware) RequestCallback(c *http.Client, r *http.Request, ctx *Context) {

}

func (p *ProxyMiddleware) Process(c *http.Client, r *http.Request, ctx *Context) {
	proxyURL := p.GetProxy()
	if len(proxyURL) > 0 {
		useProxyUrl, err := url.Parse(proxyURL)
		if err != nil {
			fmt.Println(err)
		}
		c.Transport = &http.Transport{
			Proxy: http.ProxyURL(useProxyUrl),
		}
	}
}

func (p *ProxyMiddleware) GetProxy() string {
	if len(p.List) == 0 {
		return ""
	}
	randomIndex := rand.Intn(len(p.List))
	pick := p.List[randomIndex]

	return pick
}
