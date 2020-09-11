package youcrawl

import (
	"os"
	"testing"
)

func TestOutputJsonPostProcess_Process(t *testing.T) {
	middleware := OutputJsonPostProcess{}
	middleware.StorePath = "./test_json_output.json"
	err := middleware.Process(
		&GlobalStore{
			Content: map[string]interface{}{
				"items": []map[string]interface{}{
					{
						"title": "test",
					},
					{
						"title": "test2",
					},
				},
			}})
	if err != nil {
		t.Error(err)
	}
	os.Remove(middleware.StorePath)
}
