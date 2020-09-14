package youcrawl

import "testing"

func TestMemoryGlobalStore_GetOrCreate(t *testing.T) {
	data := map[string]interface{}{
		"field": "testField",
	}
	store := MemoryGlobalStore{}
	err := store.Init()
	if err != nil {
		t.Error(err)
	}
	store.SetValue("data1", data)
	rawData := store.GetValue("data1")
	getData1 := rawData.(map[string]interface{})
	if getData1["field"] != "testField" {
		t.Error("not match")
	}

	data2 := map[string]interface{}{
		"field": "testField2",
	}

	rawData2 := store.GetOrCreate("data2", data2)
	getData2 := rawData2.(map[string]interface{})
	if getData2["field"] != "testField2" {
		t.Error("not match")
	}

}
