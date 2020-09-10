package youcrawl

import "testing"

var TestItem Item

func init() {
	TestItem = Item{Store: map[string]interface{}{}}
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
