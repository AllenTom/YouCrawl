package youcrawl

import (
	"sync"
	"testing"
)

func TestOutputJsonPostProcess_Process(t *testing.T) {
	e := NewEngine(&EngineOption{MaxRequest: 3})
	e.AddURLs("http://www.example.com")
	var wg sync.WaitGroup
	wg.Add(1)
	e.Run(&wg)
	wg.Wait()
}
