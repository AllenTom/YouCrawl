package youcrawl

type HTMLParser func(ctx *Context) error

// parse html with parser
func ParseHTML(parser HTMLParser, ctx *Context) error {
	err := parser(ctx)
	if err != nil {
		return err
	}
	return nil
}
