package youcrawl

import (
	"encoding/json"
	"io/ioutil"
)

type OutputJsonPostProcess struct {
	StorePath string
}

func (p *OutputJsonPostProcess) Process(store GlobalStore) error {
	data := store.GetValue("items")
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
