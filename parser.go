package youcrawl

import (
	"github.com/PuerkitoBio/goquery"
)

type HTMLParser func(doc *goquery.Document, ctx *Context) error

// parse html with parser
func ParseHTML(doc *goquery.Document, parser HTMLParser, ctx *Context) error {
	err := parser(doc, ctx)
	if err != nil {
		return err
	}
	return nil
}
