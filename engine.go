package youcrawl

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

var EngineLogger *logrus.Entry = logrus.WithField("scope", "engine")

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
	Pipelines   []Pipeline
	GlobalStore GlobalStore
	PostProcess []PostProcess
}

// share data in crawl process
type Context struct {
	Request     *http.Request
	Response    *http.Response
	content     map[string]interface{}
	Item        Item
	lock        *sync.Mutex
	GlobalStore *GlobalStore
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
		Pipelines:    []Pipeline{},
		Middlewares:  []Middleware{},
		Parsers:      []HTMLParser{},
		GlobalStore: GlobalStore{
			Content: map[string]interface{}{},
		},
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

// add pipelines
func (e *Engine) AddPipelines(pipelines ...Pipeline) {
	for _, pipeline := range pipelines {
		e.Pipelines = append(e.Pipelines, pipeline)
	}
}

// add middleware
func (e *Engine) UseMiddleware(middlewares ...Middleware) {
	e.Middlewares = append(e.Middlewares, middlewares...)
}

// add postprocess
func (e *Engine) AddPostProcess(postprocessList ...PostProcess) {
	e.PostProcess = append(e.PostProcess, postprocessList...)
}

// get task from pool task
func (p *RequestPool) GetTask(e *Engine) Task {
	p.Lock()
	task := p.Tasks[0]
	copy(p.Tasks, p.Tasks[1:])
	task.Context = Context{
		content:     map[string]interface{}{},
		Item:        Item{Store: map[string]interface{}{}},
		GlobalStore: &e.GlobalStore,
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
			defer func() {
				// exit run if no task
				if e.Pool.Complete() {
					// run post process
					for _, postProcess := range e.PostProcess {
						err := postProcess.Process(&e.GlobalStore)
						if err != nil {
							fmt.Println(err)
						}
					}
					stopChannel <- struct{}{}
				}
			}()
			task := e.Pool.GetTask(e)
			requestBody, err := RequestWithURL(&task, e.Middlewares...)
			if err != nil {
				EngineLogger.Error(err)
				taskChannel <- struct{}{}
				return
			}
			taskChannel <- struct{}{}
			// parse html
			// run parser one by one
			for _, parser := range e.Parsers {
				err = ParseHTML(requestBody, parser, task.Context)
				if err != nil {
					fmt.Println(err)
					return
				}
			}

			for _, pipeline := range e.Pipelines {
				err := pipeline.Process(&task.Context.Item, &e.GlobalStore)
				if err != nil {
					fmt.Println(err)
					return
				}
			}

		}()
	}

}
