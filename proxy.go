package youcrawl

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
)

var ProxyList ProxyPool

type ProxyPool struct {
	List []string
}

func init() {
	proxyList, err := ReadProxyListFile()
	if err != nil {
		fmt.Println("proxy list file cannot read,not use")
		proxyList = make([]string, 0)
	}
	ProxyList = ProxyPool{
		List: proxyList,
	}
}
func (p *ProxyPool) GetProxy() string {
	if len(p.List) == 0 {
		return ""
	}
	randomIndex := rand.Intn(len(p.List))
	pick := p.List[randomIndex]

	return pick
}
func ReadProxyListFile() ([]string, error) {
	rawList, err := ReadListFile("./proxy.txt")
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

var ProxyMiddleware Middleware = func(c *http.Client, r *http.Request, ctx *Context) {
	proxyURL := ProxyList.GetProxy()
	if len(proxyURL) > 0 {
		url, err := url.Parse(proxyURL)
		if err != nil {
			fmt.Println(err)
		}
		c.Transport = &http.Transport{
			Proxy: http.ProxyURL(url),
		}
	}
}
