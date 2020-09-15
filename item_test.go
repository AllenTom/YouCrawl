package youcrawl

import (
	"github.com/PuerkitoBio/goquery"
	"os"
	"testing"
)

var TestItem DefaultItem

func init() {
	TestItem = DefaultItem{Store: map[string]interface{}{}}
}
func TestItem_SetValue(t *testing.T) {
	TestItem.SetValue("test1", "test")
}

func TestItem_GetString(t *testing.T) {
	TestItem.SetValue("test1", "test")
	value, err := TestItem.GetString("test1")
	if err != nil {
		t.Error("get string not cause error")
	}
	if value != "test" {
		t.Error("get string value not match")
	}
	TestItem.SetValue("test1", 1)
	_, err = TestItem.GetString("test1")
	if err != TypeError {
		t.Error("must cause type error")
	}
	_, err = TestItem.GetString("test2")
	if err != KeyNotContainError {
		t.Error("must cause KeyNotContainError")
	}
}

func TestGetInt(t *testing.T) {
	TestItem.SetValue("test1", 1)
	TestItem.SetValue("test2", "1")
	_, err := TestItem.GetInt("test1")
	if err != nil {
		t.Error(err)
	}

	_, err = TestItem.GetInt("test2")
	if err != TypeError {
		t.Error("must cause type error")
	}

	_, err = TestItem.GetInt("test10")
	if err != KeyNotContainError {
		t.Error("must cause type error")
	}

}

func TestItem_GetFloat64(t *testing.T) {
	TestItem.SetValue("test1", 19.999)
	TestItem.SetValue("test2", "19.999")
	_, err := TestItem.GetFloat64("test1")
	if err != nil {
		t.Error(err)
	}
	_, err = TestItem.GetFloat64("test2")
	if err != TypeError {
		t.Error("must cause type error")
	}

	_, err = TestItem.GetFloat64("test10")
	if err != KeyNotContainError {
		t.Error("must cause type error")
	}

}

type WebItems struct {
	Label string
}
type WebInfoItem struct {
	Title  string     `json:",omitempty"`
	Items  []WebItems `json:",omitempty"`
	Status string     `json:",omitempty"`
	URL    string     `json:",omitempty"`
}

func TestItemStruct(t *testing.T) {
	e := NewEngine(&EngineOption{MaxRequest: 5})
	addTask := NewTask("http://www.bing.com", WebInfoItem{})
	e.AddTasks(&addTask)
	e.AddHTMLParser(DefaultTestParser)
	e.AddHTMLParser(func(doc *goquery.Document, ctx *Context) error {
		item := ctx.Item.(WebInfoItem)
		item.Title = doc.Find("title").Text()
		ctx.Item = item
		return nil
	})
	e.AddPipelines(&GlobalStorePipeline{})
	e.AddPostProcess(&OutputJsonPostProcess{StorePath: "./output.json"})
	e.RunAndWait()
	defer os.Remove("./output.json")
}
