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
	ID        string
	Url       string
	Context   Context
	Requested bool
	Completed bool
}

// request task pool

// youcrawl engine
type Engine struct {
	sync.Mutex
	*EngineOption
	Pool        TaskPool
	Parsers     []HTMLParser
	Middlewares []Middleware
	Pipelines   []Pipeline
	GlobalStore GlobalStore
	PostProcess []PostProcess
}

// share data in crawl process
type Context struct {
	sync.Mutex
	Request     *http.Request
	Response    *http.Response
	content     map[string]interface{}
	Item        Item
	GlobalStore GlobalStore
	Pool        *RequestPool
}

// init engine config
type EngineOption struct {
	MaxRequest int
}

// init new engine
func NewEngine(option *EngineOption) *Engine {
	globalStore := MemoryGlobalStore{}
	err := globalStore.Init()
	if err != nil {
		logrus.Fatal("init global store failed")
	}
	newEngine := &Engine{
		Pool: &RequestPool{
			Tasks:    []Task{},
			DoneChan: make(chan struct{}),
		},
		EngineOption: option,
		Pipelines:    []Pipeline{},
		Middlewares:  []Middleware{},
		Parsers:      []HTMLParser{},
		GlobalStore:  &globalStore,
	}

	return newEngine
}

// add url to crawl
// unsafe operation,engine must not in running status
//
// in engine running ,use RequestPool.AddURLs method
func (e *Engine) AddURLs(urls ...string) {
	e.Pool.AddURLs(urls...)
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

func CrawlProcess(taskChannel chan struct{}, e *Engine, task *Task) {
	defer e.Pool.OnTaskDone(task)
	requestBody, err := RequestWithURL(task, e.Middlewares...)
	if err != nil {
		EngineLogger.Info(err)
		taskChannel <- struct{}{}
		return
	}
	taskChannel <- struct{}{}
	// parse html
	// run parser one by one
	for _, parser := range e.Parsers {
		err = ParseHTML(requestBody, parser, &task.Context)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}

	for _, pipeline := range e.Pipelines {
		err := pipeline.Process(&task.Context.Item, e.GlobalStore)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}

}

// run crawl engine
func (e *Engine) Run(wg *sync.WaitGroup) {
	defer func() {
		EngineLogger.Info("all done ,send stop signal")
		wg.Done()
	}()
	taskChannel := make(chan struct{}, e.MaxRequest)
	for idx := 0; idx < e.MaxRequest; idx++ {
		taskChannel <- struct{}{}
	}
Loop:
	for {
		<-taskChannel
		select {
		case task := <-e.Pool.GetOneTask(e):

			go CrawlProcess(taskChannel, e, task)
		case <-e.Pool.GetDoneChan():
			break Loop
		}
	}

	EngineLogger.Info("into post process")
	for _, postProcess := range e.PostProcess {
		err := postProcess.Process(e.GlobalStore)
		if err != nil {
			fmt.Println(err)
		}
	}
	return
}
