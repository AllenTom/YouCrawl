package youcrawl

type PostProcess interface {
	Process(store GlobalStore) error
}
