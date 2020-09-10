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
	Url       string
	Context   Context
	Requested bool
	Completed bool
}

// request task pool
type RequestPool struct {
	Tasks         []Task
	Total         int
	CompleteCount int
	NextTask      *Task
	GetTaskChan   chan *Task
	DoneChan      chan struct{}
	sync.Mutex
}

func (p *RequestPool) GetOneTask(e *Engine) <-chan *Task {
	taskChan := make(chan *Task)
	go func(callbackChan chan *Task) {
		p.Lock()
		defer p.Unlock()
		for idx := range p.Tasks {
			current := &p.Tasks[idx]
			current.Init(p, &e.GlobalStore)
			if !current.Completed && !current.Requested {
				current.Requested = true
				callbackChan <- current
				return
			}
		}
		// no more request,wait for new
		p.GetTaskChan = callbackChan
	}(taskChan)
	return taskChan
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
	Pool        *RequestPool
}

// init engine config
type EngineOption struct {
	MaxRequest int
}

// init new engine
func NewEngine(option *EngineOption) *Engine {
	newEngine := &Engine{
		Pool: &RequestPool{
			Tasks:    []Task{},
			DoneChan: make(chan struct{}),
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

func (t *Task) Init(p *RequestPool, g *GlobalStore) {
	t.Context = Context{
		Pool:        p,
		content:     map[string]interface{}{},
		Item:        Item{Store: map[string]interface{}{}},
		GlobalStore: g,
	}
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
func (p *RequestPool) Resume(e *Engine) {
	p.Lock()
	defer p.Unlock()
	// no hang up,skip

	// find not run task
	for idx := range e.Pool.Tasks {
		targetTask := &e.Pool.Tasks[idx]
		if !targetTask.Requested && !targetTask.Completed {
			targetTask.Init(e.Pool, &e.GlobalStore)
			e.Pool.GetTaskChan <- targetTask
			return
		}
	}

	// find not complete
	for idx := range e.Pool.Tasks {
		targetTask := &e.Pool.Tasks[idx]
		if !targetTask.Completed {
			return
		}
	}

	// no new request to resume
	EngineLogger.Info("no more task to resume , go to done!")
	p.DoneChan <- struct{}{}

}

func CrawlProcess(taskChannel chan struct{}, e *Engine, task *Task) {
	defer func() {
		// exit run if no task
		task.Completed = true
		// mark it done

		e.Pool.Resume(e)
		// current is last,no more task,switch to done
	}()
	requestBody, err := RequestWithURL(task, e.Middlewares...)

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
			continue
		}
	}

	for _, pipeline := range e.Pipelines {
		err := pipeline.Process(&task.Context.Item, &e.GlobalStore)
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
			go func() {
				e.Pool.Lock()
				e.Pool.GetTaskChan = nil
				defer e.Pool.Unlock()
				go CrawlProcess(taskChannel, e, task)
			}()
		case <-e.Pool.DoneChan:
			break Loop
		}
	}

	EngineLogger.Info("into post process")
	for _, postProcess := range e.PostProcess {
		err := postProcess.Process(&e.GlobalStore)
		if err != nil {
			fmt.Println(err)
		}
	}
	return
}
