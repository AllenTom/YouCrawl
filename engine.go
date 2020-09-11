package youcrawl

import (
	"fmt"
	"github.com/rs/xid"
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
type RequestPool struct {
	Tasks         []Task
	Total         int
	CompleteCount int
	NextTask      *Task
	GetTaskChan   chan *Task
	DoneChan      chan struct{}
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
	sync.Mutex
	Request     *http.Request
	Response    *http.Response
	content     map[string]interface{}
	Item        Item
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

// add task to task pool
func (p *RequestPool) AddURLs(urls ...string) {
	EngineLogger.Info(fmt.Sprintf("append new url with len = %d", len(urls)))
	p.Lock()
	defer p.Unlock()
	p.Total += len(urls)
	for _, url := range urls {
		p.Tasks = append(p.Tasks, Task{
			ID:  xid.New().String(),
			Url: url,
			Context: Context{
				content: map[string]interface{}{},
				Item: Item{
					Store: map[string]interface{}{},
				},
			},
		})
	}

	// suspend task requirement exist,resume
	// see also `RequestPool.GetOneTask` method
	if p.GetTaskChan != nil {
		resumeTask := p.GetUnRequestedTask()
		if resumeTask != nil {
			resumeTask.Context.Pool = p
			p.GetTaskChan <- resumeTask
			resumeTask.Requested = true
			p.GetTaskChan = nil
		}
	}
}

func (p *RequestPool) GetOneTask(e *Engine) <-chan *Task {
	taskChan := make(chan *Task)
	go func(callbackChan chan *Task) {
		p.Lock()
		defer p.Unlock()
		unRequestedTask := p.GetUnRequestedTask()
		if unRequestedTask != nil {
			unRequestedTask.Context.Pool = p
			unRequestedTask.Requested = true
			unRequestedTask.Context.GlobalStore = &e.GlobalStore
			callbackChan <- unRequestedTask
			return
		}
		// no more request,suspend task
		EngineLogger.Info("suspend get task ")
		p.GetTaskChan = callbackChan
	}(taskChan)
	return taskChan
}

// find unreauested task
func (p *RequestPool) GetUnRequestedTask() (target *Task) {
	for idx := range p.Tasks {
		iterTask := &p.Tasks[idx]
		if !iterTask.Requested {
			target = iterTask
			return
		}
	}
	return nil
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

// get task from pool task
func (p *RequestPool) OnTaskDone(task *Task) {
	p.Lock()
	defer p.Unlock()
	// set true

	for idx := range p.Tasks {
		if task.ID == p.Tasks[idx].ID {
			p.Tasks[idx].Completed = true
		}
	}

	// check weather all task are complete
	for idx := range p.Tasks {
		if !p.Tasks[idx].Requested || !p.Tasks[idx].Completed {
			return
		}
	}

	// no new request
	EngineLogger.Info("no more task to resume , go to done!")
	p.DoneChan <- struct{}{}

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

			go CrawlProcess(taskChannel, e, task)
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
