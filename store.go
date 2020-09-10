package youcrawl

import "sync"

//store engine global
type GlobalStore struct {
	sync.Mutex
	Content map[string]interface{}
}
