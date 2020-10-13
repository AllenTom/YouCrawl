package youcrawl

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

type Plugin interface {
	Run(e *Engine)
}
const (
	// total
	STATUS_KEY_TOTAL = "status.total"
	// unrequested count
	STATUS_KEY_UNREQUESTED = "status.unrequested"
	// complete count
	STATUS_KEY_COMPLETE = "status.complete"
	// speed
	STATUS_KEY_SPEED = "status.speed"
)

// log engine status plugin
type StatusOutputPlugin struct {
	// disable log output
	LogOutput bool
}

func (p *StatusOutputPlugin) Run(e *Engine) {
	lastComplete := 0
	for  {
		total, _ := e.Pool.GetTotal()
		unrequested, _ := e.Pool.GetUnRequestCount()
		complete, _ := e.Pool.GetCompleteCount()
		speed := complete - lastComplete
		if p.LogOutput {
			logrus.WithField("scope","status-report").Info(fmt.Sprintf(
				"total: %d,unrequest: %d,complete: %d,speed: %d/s",
				total,unrequested,complete,speed,
			))
		}
		e.GlobalStore.SetValue(STATUS_KEY_TOTAL,total)
		e.GlobalStore.SetValue(STATUS_KEY_UNREQUESTED,unrequested)
		e.GlobalStore.SetValue(STATUS_KEY_COMPLETE,complete)
		e.GlobalStore.SetValue(STATUS_KEY_SPEED,speed)
		lastComplete = complete
		<-time.After(1 * time.Second)
	}

}
