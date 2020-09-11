package youcrawl

import (
	"sync"
)

//store engine global
type GlobalStore interface {
	Init() error
	SetValue(key string, value interface{})
	GetValue(key string) interface{}
}
type MemoryGlobalStore struct {
	sync.Mutex
	Content map[string]interface{}
}

func (s *MemoryGlobalStore) SetValue(key string, value interface{}) {
	s.Lock()
	defer s.Unlock()

}

func (s *MemoryGlobalStore) GetValue(key string) interface{} {
	target, exist := s.Content[key]
	if exist {
		return target
	} else {
		return nil
	}
}
func (s *MemoryGlobalStore) Init() error {
	s.Lock()
	defer s.Unlock()

	s.Content = map[string]interface{}{}
	return nil
}

// global store pipeline
// save current item to global items
type GlobalStorePipeline struct {
}

func (g *GlobalStorePipeline) Process(item *Item, store GlobalStore) error {
	rawItems := store.GetValue("items")
	if rawItems != nil {
		items := rawItems.([]map[string]interface{})
		items = append(items, item.Store)
		store.SetValue("items", items)
	}

	return nil
}
