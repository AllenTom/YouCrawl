package youcrawl

import (
	"sync"
)

type Request struct {
	Url string
}
type RequestPool struct {
	Tasks         []Request
	Total         int
	CompleteCount int
	sync.Mutex
}
type Engine struct {
	sync.Mutex
	*EngineOption
	Pool             *RequestPool
	RunningTaskCount int
}
type EngineOption struct {
	MaxRequest int
}

func NewEngine(option *EngineOption) *Engine {
	newEngine := &Engine{
		RunningTaskCount: 0,
		Pool: &RequestPool{
			Tasks:         []Request{},
			CompleteCount: 0,
		},
		EngineOption: option,
	}

	return newEngine
}

func (e *Engine) AddURLs(urls ...string) {
	e.Pool.Total += len(urls)
	for _, url := range urls {
		e.Pool.Tasks = append(e.Pool.Tasks, Request{Url: url})
	}
}

func (p *RequestPool) GetTask() Request {
	p.Lock()
	var task Request
	task, p.Tasks = p.Tasks[0], p.Tasks[1:]
	defer p.Unlock()
	return task
}
func (p *RequestPool) Complete() bool {
	p.Lock()
	defer p.Unlock()
	p.CompleteCount += 1
	if p.Total == p.CompleteCount {
		return true
	}
	return false
}

func (e *Engine) Run(stopChannel chan<- struct{}) {
	taskChannel := make(chan struct{}, e.MaxRequest)
	for idx := 0; idx < e.MaxRequest; idx++ {
		taskChannel <- struct{}{}
	}
	for idx := 0; idx < len(e.Pool.Tasks); idx++ {
		go func() {
			<-taskChannel
			task := e.Pool.GetTask()
			err := RequestWithURL(task.Url)
			if err != nil {
				taskChannel <- struct{}{}
				return
			}
			taskChannel <- struct{}{}

			// exit run if no task
			if e.Pool.Complete() {
				stopChannel <- struct{}{}
			}
		}()
	}

}
