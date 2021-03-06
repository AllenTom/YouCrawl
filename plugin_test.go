package youcrawl

import (
	"fmt"
	"testing"
	"time"
)

type HeartbeatPlugin struct {

}

func (p *HeartbeatPlugin) Run(e *Engine) {
	for {
		<-time.After(1*time.Second)
		fmt.Println("heartbeat log of engine")
	}
}

func TestPlugins(t *testing.T) {
	e := NewEngine(&EngineOption{
		MaxRequest: 3,
		Daemon:     true,
	})
	plugin := HeartbeatPlugin{}
	e.AddPlugins(&plugin)
	e.AddURLs("http://www.example.com")
	go e.RunAndWait()
	<-time.After(2*time.Second)
	e.StopPoolChan <- struct{}{}
}

func TestStatusOutputPlugin_Run(t *testing.T) {
	e := NewEngine(&EngineOption{
		MaxRequest: 2,
		Daemon:     false,
	})
	plugin := StatusOutputPlugin{
		LogOutput: true,
	}
	e.AddPlugins(&plugin)
	e.AddURLs("http://www.example.com")
	go func() {
		for {
			fmt.Println(e.GlobalStore.GetValue(STATUS_KEY_SPEED))
			<-time.After(1 * time.Second)
		}
	}()
	e.RunAndWait()
}