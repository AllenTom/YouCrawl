package youcrawl

import (
	"encoding/json"
	"io/ioutil"
)

type OutputJsonPostProcess struct {
	StorePath string
}

func (p *OutputJsonPostProcess) Process(store *GlobalStore) error {
	file, _ := json.MarshalIndent(store.Content["items"], "", " ")
	_ = ioutil.WriteFile(p.StorePath, file, 0644)
	return nil
}
