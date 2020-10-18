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
	SetPrevent(isPrevent bool)
	GetTotal() (int, error)
	GetUnRequestCount() (int, error)
	GetCompleteCount() (int, error)
}

type RequestPool struct {
	Tasks         []Task
	Total         int
	CompleteCount int
	NextTask      *Task
	GetTaskChan   chan *Task
	DoneChan      chan struct{}
	CompleteChan  chan *Task
	PreventStop   bool
	Store         GlobalStore
	sync.RWMutex
}

func (p *RequestPool) GetCompleteCount() (int, error) {
	p.Lock()
	defer p.Unlock()
	count := 0
	for _, task := range p.Tasks {
		if task.Completed {
			count += 1
		}
	}
	return count, nil
}

func (p *RequestPool) GetUnRequestCount() (int, error) {
	p.Lock()
	defer p.Unlock()
	count := 0
	for _, task := range p.Tasks {
		if !task.Requested {
			count += 1
		}
	}
	return count, nil
}

func (p *RequestPool) GetTotal() (int, error) {
	p.Lock()
	defer p.Unlock()
	return p.Total, nil
}

func (p *RequestPool) SetPrevent(isPrevent bool) {
	p.Lock()
	defer p.Unlock()
	p.PreventStop = isPrevent

	if !isPrevent {
		for idx := range p.Tasks {
			if !p.Tasks[idx].Requested || !p.Tasks[idx].Completed {
				// not done task exist, not kill
				return
			}
		}
		EngineLogger.Info("no more task to resume , go to done!")
		p.DoneChan <- struct{}{}
	}

}

type RequestPoolOption struct {
	UseCookie   bool
	PreventStop bool
}

func NewRequestPool(option RequestPoolOption, store GlobalStore) *RequestPool {
	pool := &RequestPool{
		Tasks:         []Task{},
		DoneChan:      make(chan struct{}),
		Total:         0,
		CompleteCount: 0,
		Store:         store,
		PreventStop:   option.PreventStop,
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
	if p.GetTaskChan != nil {
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
	if p.GetTaskChan != nil {
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
		unRequestedTask := p.GetUnRequestedTask()
		if unRequestedTask != nil {
			unRequestedTask = p.initTask(unRequestedTask)
			unRequestedTask.Requested = true
			callbackChan <- unRequestedTask
		}
		waitChannel := make(chan *Task)
		p.GetTaskChan = waitChannel
		p.Unlock()
		callbackChan <- <-waitChannel
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

	// no new request
	if !p.PreventStop {
		EngineLogger.Info("no more task to resume , go to done!")
		p.DoneChan <- struct{}{}
	}
}

func (p *RequestPool) GetDoneChan() chan struct{} {
	return p.DoneChan
}

func (p *RequestPool) Close() {
	p.Lock()
	defer p.Unlock()
	p.DoneChan <- struct{}{}
}
