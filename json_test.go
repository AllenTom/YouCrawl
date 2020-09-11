package youcrawl

import (
	"os"
	"testing"
)

func TestOutputJsonPostProcess_Process(t *testing.T) {
	gb := MemoryGlobalStore{}
	err := gb.Init()
	if err != nil {
		t.Error(err)
	}
	gb.SetValue("items", []map[string]interface{}{
		{
			"title": "test",
		},
		{
			"title": "test2",
		},
	})
	middleware := OutputJsonPostProcess{}
	middleware.StorePath = "./test_json_output.json"
	err = middleware.Process(&gb)
	if err != nil {
		t.Error(err)
	}
	os.Remove(middleware.StorePath)
}
