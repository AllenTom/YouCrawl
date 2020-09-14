package youcrawl

import (
	"sync"
)

//store engine global
type GlobalStore interface {
	Init() error
	SetValue(key string, value interface{})
	GetValue(key string) interface{}
	GetOrCreate(key string, value interface{}) interface{}
}
type MemoryGlobalStore struct {
	sync.Map
	Content map[string]interface{}
}

func (s *MemoryGlobalStore) GetOrCreate(key string, value interface{}) interface{} {
	result, _ := s.LoadOrStore(key, value)
	return result
}

func (s *MemoryGlobalStore) SetValue(key string, value interface{}) {
	s.Store(key, value)
}

func (s *MemoryGlobalStore) GetValue(key string) interface{} {
	target, exist := s.Load(key)
	if exist {
		return target
	} else {
		return nil
	}
}
func (s *MemoryGlobalStore) Init() error {
	return nil
}

// global store pipeline
// save current item to global items
type GlobalStorePipeline struct {
}

func (g *GlobalStorePipeline) Process(item *Item, store GlobalStore) error {
	rawItems := store.GetValue("items")
	if rawItems == nil {
		rawItems = make([]map[string]interface{}, 0)
	}

	items := rawItems.([]map[string]interface{})
	items = append(items, item.Store)
	store.SetValue("items", items)
	return nil
}
