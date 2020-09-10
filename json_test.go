package youcrawl

import "testing"

func TestOutputJsonPostProcess_Process(t *testing.T) {
	e := NewEngine(&EngineOption{MaxRequest: 3})
	e.AddURLs("http://www.example.com")
	stopChannel := make(chan struct{})
	e.Run(stopChannel)
	<-stopChannel
}
