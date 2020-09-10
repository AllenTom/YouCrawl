package youcrawl

import (
	"sync"
)

//store engine global
type GlobalStore struct {
	sync.Mutex
	Content map[string]interface{}
}

// global store pipeline
// save current item to global items
type GlobalStorePipeline struct {
}

func (g *GlobalStorePipeline) Process(item *Item, store *GlobalStore) error {
	store.Lock()
	defer store.Unlock()
	if store.Content["items"] == nil {
		store.Content["items"] = make([]map[string]interface{}, 0)
	}
	rawItems := store.Content["items"]
	items := rawItems.([]map[string]interface{})
	store.Content["items"] = append(items, item.Store)
	return nil
}
