package youcrawl

import (
	"bufio"
	"fmt"
	"math/rand"
	"net/http"
	"os"
)

var UserAgents UserAgentPool

// random user agent middleware
type UserAgentPool struct {
	List []string
}

func init() {
	agents, err := ReadUserAgentListFile()
	if err != nil {
		fmt.Println("read ua file fail,no ua will be used")
		agents = make([]string, 0)
	}
	UserAgents = UserAgentPool{
		List: agents,
	}
}
func (p *UserAgentPool) GetUserAgent() string {
	randomIndex := rand.Intn(len(p.List))
	pick := p.List[randomIndex]
	fmt.Println(fmt.Sprintf("pick ua %d", randomIndex))
	return pick
}
func ReadUserAgentListFile() ([]string, error) {
	file, err := os.Open("./ua.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	uaList := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		uaList = append(uaList, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return uaList, err
}

var UserAgentMiddleware Middleware = func(r *http.Request, ctx Context) {
	r.Header.Add("User-Agent", UserAgents.GetUserAgent())
}
