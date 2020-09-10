package youcrawl

type Pipeline interface {
	Process(item *Item, store *GlobalStore) error
}
