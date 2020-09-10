package youcrawl

type Pipeline interface {
	Process(item *Item) error
}
