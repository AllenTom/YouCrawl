package youcrawl

import (
	"github.com/PuerkitoBio/goquery"
	"io"
)

type HTMLParser func(doc *goquery.Document,ctx Context)

// parse html with parser
func ParseHTML(reader io.Reader,parser HTMLParser,ctx Context) error{
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return err
	}
	parser(doc,ctx)
	return nil
}