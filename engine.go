package youcrawl

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/cookiejar"
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
	// dispatch task
	Pool          TaskPool
	Parsers       []HTMLParser
	Middlewares   []Middleware
	Pipelines     []Pipeline
	GlobalStore   GlobalStore
	PostProcess   []PostProcess
	InterruptChan chan struct{}
}

// share data in crawl process
type Context struct {
	sync.Mutex
	Request     *http.Request
	Response    *http.Response
	Item        interface{}
	GlobalStore GlobalStore
	Pool        TaskPool
	Cookie      *cookiejar.Jar
}

// init engine config
type EngineOption struct {
	// max running in same time
	MaxRequest int
}

// init new engine
func NewEngine(option *EngineOption) *Engine {
	globalStore := &MemoryGlobalStore{}
	err := globalStore.Init()
	if err != nil {
		logrus.Fatal("init global store failed")
	}
	pool := NewRequestPool(RequestPoolOption{}, globalStore)

	newEngine := &Engine{
		Pool:          pool,
		EngineOption:  option,
		Pipelines:     []Pipeline{},
		Middlewares:   []Middleware{},
		Parsers:       []HTMLParser{},
		GlobalStore:   globalStore,
		InterruptChan: make(chan struct{}),
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

// add task to crawl
// unsafe operation,engine must not in running status
//
// in engine running ,use RequestPool.AddURLs method
func (e *Engine) AddTasks(tasks ...*Task) {
	e.Pool.AddTasks(tasks...)
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
	doc, err := goquery.NewDocumentFromReader(requestBody)
	if err != nil {
		EngineLogger.Error(err)
	}
	// run parser one by one

	for _, parser := range e.Parsers {
		err = ParseHTML(doc, parser, &task.Context)
		if err != nil {
			EngineLogger.Error(err)
			continue
		}
	}

	for _, pipeline := range e.Pipelines {
		err := pipeline.Process(task.Context.Item, e.GlobalStore)
		if err != nil {
			EngineLogger.Error(err)
			continue
		}
	}

}

// run and wait it done
func (e *Engine) RunAndWait() {
	var wg sync.WaitGroup
	wg.Add(1)
	e.Run(&wg)
	wg.Wait()
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

	// run interrupt chan
	go func() {
		select {
		case <-e.InterruptChan:
			e.Pool.Close()
		}
	}()
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
