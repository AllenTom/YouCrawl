package youcrawl

import (
	"fmt"
	"github.com/rs/xid"
	"sync"
)

type TaskPool interface {
	AddURLs(urls ...string)
	GetOneTask(e *Engine) <-chan *Task
	GetUnRequestedTask() (target *Task)
	OnTaskDone(task *Task)
	GetDoneChan() chan struct{}
}

type RequestPool struct {
	Tasks         []Task
	Total         int
	CompleteCount int
	NextTask      *Task
	GetTaskChan   chan *Task
	DoneChan      chan struct{}
	sync.Mutex
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

func (p *RequestPool) GetDoneChan() chan struct{} {
	return p.DoneChan
}
