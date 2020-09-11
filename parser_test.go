package youcrawl

import (
	"testing"
)

func TestParser(t *testing.T) {
	task := &Task{
		Url: "http://www.example.com",
		Context: Context{
			content: map[string]interface{}{},
		},
	}
	reader, err := RequestWithURL(task)
	if err != nil {
		t.Error(err)
	}
	err = ParseHTML(reader, DefaultTestParser, &task.Context)
	if err != nil {
		t.Error(err)
	}
}
