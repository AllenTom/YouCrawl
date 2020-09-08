package youcrawl

import "testing"

func TestRequestWithURL(t *testing.T) {
	err := RequestWithURL("https://www.zhihu.com")
	if err != nil {
		t.Error(err)
	}
}

func TestEngine_Run(t *testing.T) {
	e := NewEngine(&EngineOption{MaxRequest: 3})
	e.AddURLs("https://www.zhihu.com")
	stopChannel := make(chan struct{})
	e.Run(stopChannel)
	<-stopChannel
}
