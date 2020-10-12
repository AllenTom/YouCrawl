package youcrawl

type Plugin interface {
	Run(e *Engine)
}
