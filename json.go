package youcrawl

import (
	"encoding/json"
	"io/ioutil"
)

type OutputJsonPostProcess struct {
	StorePath string
	GetData   func(store GlobalStore) interface{}
}

func (p *OutputJsonPostProcess) Process(store GlobalStore) error {
	data := store.GetValue("items")
	if p.GetData != nil {
		data = p.GetData(store)
	}
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(p.StorePath, file, 0644)
	if err != nil {
		return err
	}
	return nil
}
