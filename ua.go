package youcrawl

import (
	"fmt"
	"math/rand"
	"net/http"
)

var UserAgents UserAgentPool

// random user agent middleware
type UserAgentPool struct {
	List []string
}

func init() {
	agents, err := ReadListFile("./ua.txt")
	if err != nil {
		fmt.Println("read ua file fail,no ua will be used")
		agents = make([]string, 0)
	}
	UserAgents = UserAgentPool{
		List: agents,
	}
}
func (p *UserAgentPool) GetUserAgent() string {
	if len(p.List) == 0 {
		return ""
	}
	randomIndex := rand.Intn(len(p.List))
	pick := p.List[randomIndex]
	return pick
}

type UserAgentMiddleware struct{}

func (u *UserAgentMiddleware) RequestCallback(c *http.Client, r *http.Request, ctx *Context) {

}

func (u *UserAgentMiddleware) Process(c *http.Client, r *http.Request, ctx *Context) {
	ua := UserAgents.GetUserAgent()
	if len(ua) > 0 {
		r.Header.Add("User-Agent", ua)
	}
}
