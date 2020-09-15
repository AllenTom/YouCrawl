package youcrawl

import (
	"fmt"
	"github.com/rs/xid"
	"sync"
)

type TaskPool interface {
	AddURLs(urls ...string)
	AddTasks(task ...*Task)
	GetOneTask(e *Engine) <-chan *Task
	GetUnRequestedTask() (target *Task)
	OnTaskDone(task *Task)
	GetDoneChan() chan struct{}
	Close()
}

type RequestPool struct {
	Tasks         []Task
	Total         int
	CompleteCount int
	NextTask      *Task
	GetTaskChan   chan *Task
	DoneChan      chan struct{}
	CloseFlag     int64
	CompleteChan  chan *Task
	UseCookie     bool
	Store         GlobalStore
	sync.Mutex
}
type RequestPoolOption struct {
	UseCookie bool
}

func NewRequestPool(option RequestPoolOption, store GlobalStore) *RequestPool {
	pool := &RequestPool{
		Tasks:         []Task{},
		DoneChan:      make(chan struct{}),
		CloseFlag:     0,
		Total:         0,
		CompleteCount: 0,
		UseCookie:     option.UseCookie,
		Store:         store,
	}
	return pool
}

func NewTask(url string, item interface{}) Task {
	return Task{
		ID:  xid.New().String(),
		Url: url,
		Context: Context{
			Item: item,
		},
	}
}

// add task
func (p *RequestPool) AddTasks(tasks ...*Task) {
	EngineLogger.Info(fmt.Sprintf("append new url with len = %d", len(tasks)))
	p.Lock()
	defer p.Unlock()
	p.Total += len(tasks)
	for _, addTask := range tasks {
		item := addTask.Context.Item
		if item == nil {
			item = DefaultItem{
				Store: map[string]interface{}{},
			}
		}
		p.Tasks = append(p.Tasks, Task{
			ID:  addTask.ID,
			Url: addTask.Url,
			Context: Context{
				Item: addTask.Context.Item,
			},
		})
	}

	// suspend task requirement exist,resume
	// see also `RequestPool.GetOneTask` method
	if p.GetTaskChan != nil && p.CloseFlag == 0 {
		resumeTask := p.GetUnRequestedTask()
		if resumeTask != nil {
			resumeTask = p.initTask(resumeTask)
			resumeTask.Requested = true
			p.GetTaskChan <- resumeTask
			p.GetTaskChan = nil
		}
	}
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
				Item: DefaultItem{
					Store: map[string]interface{}{},
				},
			},
		})
	}

	// suspend task requirement exist,resume
	// see also `RequestPool.GetOneTask` method
	if p.GetTaskChan != nil && p.CloseFlag == 0 {
		resumeTask := p.GetUnRequestedTask()
		if resumeTask != nil {
			resumeTask = p.initTask(resumeTask)
			resumeTask.Requested = true
			p.GetTaskChan <- resumeTask
			p.GetTaskChan = nil
		}
	}
}
func (p *RequestPool) initTask(task *Task) *Task {
	task.Context.Pool = p
	task.ID = xid.New().String()
	task.Context.GlobalStore = p.Store
	return task
}
func (p *RequestPool) GetOneTask(e *Engine) <-chan *Task {
	taskChan := make(chan *Task)
	go func(callbackChan chan *Task) {
		p.Lock()
		defer p.Unlock()
		if p.CloseFlag == 1 {
			return
		}
		unRequestedTask := p.GetUnRequestedTask()
		if unRequestedTask != nil {
			unRequestedTask = p.initTask(unRequestedTask)
			unRequestedTask.Requested = true
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
	if p.CloseFlag == 0 {
		for idx := range p.Tasks {
			if !p.Tasks[idx].Requested || !p.Tasks[idx].Completed {
				return
			}
		}
	}

	// no new request
	EngineLogger.Info("no more task to resume , go to done!")
	p.DoneChan <- struct{}{}

}

func (p *RequestPool) GetDoneChan() chan struct{} {
	return p.DoneChan
}

func (p *RequestPool) Close() {
	p.Lock()
	defer p.Unlock()
	p.CloseFlag = 1
}
