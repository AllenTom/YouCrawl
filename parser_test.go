package youcrawl

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"testing"
)

func TestParser(t *testing.T) {
	task := &Task{
		Url:     "http://www.example.com",
		Context: Context{},
	}
	reader, err := RequestWithURL(task)
	if err != nil {
		t.Error(err)
	}
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		EngineLogger.Error(err)
	}
	err = ParseHTML(doc, DefaultTestParser, &task.Context)
	if err != nil {
		t.Error(err)
	}
}

func TestParserOnError(t *testing.T) {
	task := &Task{
		Url:     "https://api.github.com/",
		Context: Context{},
	}
	err := ParseHTML(nil, func(doc *goquery.Document, ctx *Context) error {
		return errors.New("test error")
	}, &task.Context)
	if err == nil {
		t.Error("must cause error")
	}
}
