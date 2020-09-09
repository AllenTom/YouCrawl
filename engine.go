package youcrawl

import (
	"fmt"
	"net/http"
	"sync"
)

// tracking request task
type Task struct {
	Url     string
	Context Context
}

// request task pool
type RequestPool struct {
	Tasks         []Task
	Total         int
	CompleteCount int
	sync.Mutex
}

// youcrawl engine
type Engine struct {
	sync.Mutex
	*EngineOption
	Pool        *RequestPool
	Parsers     []HTMLParser
	Middlewares []Middleware
}

// share data in crawl process
type Context struct {
	Request  *http.Request
	Response *http.Response
	content  map[string]interface{}
	lock     *sync.Mutex
}

// init engine config
type EngineOption struct {
	MaxRequest int
}

// init new engine
func NewEngine(option *EngineOption) *Engine {
	newEngine := &Engine{
		Pool: &RequestPool{
			Tasks:         []Task{},
			CompleteCount: 0,
		},
		EngineOption: option,
	}

	return newEngine
}

// add url to crawl
func (e *Engine) AddURLs(urls ...string) {
	e.Pool.Total += len(urls)
	for _, url := range urls {
		e.Pool.Tasks = append(e.Pool.Tasks, Task{Url: url})
	}
}

// add parse
func (e *Engine) AddHTMLParser(parsers ...HTMLParser) {
	for _, htmlParser := range parsers {
		e.Parsers = append(e.Parsers, htmlParser)
	}

}

// add middleware
func (e *Engine) UseMiddleware(middlewares ...Middleware) {
	e.Middlewares = append(e.Middlewares, middlewares...)
}

// get task from pool task
func (p *RequestPool) GetTask() Task {
	p.Lock()
	task := p.Tasks[0]
	copy(p.Tasks, p.Tasks[1:])
	task.Context = Context{
		content: map[string]interface{}{},
	}
	defer p.Unlock()
	return task
}

// complete task
func (p *RequestPool) Complete() bool {
	p.Lock()
	defer p.Unlock()
	p.CompleteCount += 1
	if p.Total == p.CompleteCount {
		return true
	}
	return false
}

// run crawl engine
func (e *Engine) Run(stopChannel chan<- struct{}) {
	taskChannel := make(chan struct{}, e.MaxRequest)
	for idx := 0; idx < e.MaxRequest; idx++ {
		taskChannel <- struct{}{}
	}
	for idx := 0; idx < len(e.Pool.Tasks); idx++ {
		go func() {
			<-taskChannel
			task := e.Pool.GetTask()
			requestBody, err := RequestWithURL(&task, e.Middlewares...)
			if err != nil {
				taskChannel <- struct{}{}
				return
			}
			taskChannel <- struct{}{}
			// parse html
			// run parser one by one
			for _, parser := range e.Parsers {
				var parseWg sync.WaitGroup
				parseWg.Add(1)
				go func(wg *sync.WaitGroup) {
					defer parseWg.Done()
					err = ParseHTML(requestBody, parser, task.Context)
					if err != nil {
						fmt.Print(err)
					}
				}(&parseWg)
				parseWg.Wait()
			}

			// exit run if no task
			if e.Pool.Complete() {
				stopChannel <- struct{}{}
			}
		}()
	}

}
