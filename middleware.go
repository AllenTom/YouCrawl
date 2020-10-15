package youcrawl

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"
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

type UserAgentMiddleware struct {
	List []string
}

type UserAgentMiddlewareOption struct {
	// set user agent list,
	// if both UserAgentList and UserAgentFilePath are provided,combine tow list
	UserAgentList []string
	// read useragent from file,use `./ua.txt` by default,
	// if both UserAgentList and UserAgentFilePath are provided,combine tow list
	UserAgentFilePath string
}

func NewUserAgentMiddleware(option UserAgentMiddlewareOption) (*UserAgentMiddleware, error) {
	UserAgentList := option.UserAgentList
	if UserAgentList == nil {
		UserAgentList = make([]string, 0)
	}

	list, err := readUserAgentListFile(option.UserAgentFilePath)
	if err != nil {
		return nil, err
	}
	UserAgentList = append(UserAgentList, list...)
	middleware := &UserAgentMiddleware{List: UserAgentList}
	return middleware, nil
}

func readUserAgentListFile(UserAgentPath string) ([]string, error) {
	if len(UserAgentPath) == 0 {
		UserAgentPath = "./ua.txt"
	}
	rawList, err := ReadListFile(UserAgentPath)
	if err != nil {
		return nil, err
	}
	return rawList, nil
}

func (p *UserAgentMiddleware) RequestCallback(c *http.Client, r *http.Request, ctx *Context) {

}

func (p *UserAgentMiddleware) Process(c *http.Client, r *http.Request, ctx *Context) {
	uaString := p.GetUserAgent()
	if len(uaString) > 0 {
		r.Header.Add("User-Agent", uaString)
	}
}

func (p *UserAgentMiddleware) GetUserAgent() string {
	if len(p.List) == 0 {
		return ""
	}
	randomIndex := rand.Intn(len(p.List))
	pick := p.List[randomIndex]

	return pick
}

type DelayMiddleware struct {
	Min int
	Max int
	Fixed int
}

func (d *DelayMiddleware) Process(_ *http.Client, r *http.Request, _ *Context) {
	if d.Fixed != 0 {
		<- time.After(time.Duration(d.Fixed) * time.Second)
	}else if d.Min < d.Max{
		randomValue := RandomIntRangeWithStringSeed(d.Min,d.Max,r.URL.String())
		<- time.After(time.Duration(randomValue) * time.Second)
	}
}

func (d *DelayMiddleware) RequestCallback(_ *http.Client, _ *http.Request, _ *Context) {

}