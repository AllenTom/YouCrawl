# Parser
Parser is a component that processes HTML documents and extracts information from the page.

Parser use: [goquery](https://github.com/PuerkitoBio/goquery) library

## Custom Parser
Use `e.AddHTMLParser(yourParser)` to add a processing function.

## Item
The `Item` is placed in the `Context`, and the information obtained in the Parser is saved and used by the next component (Pipeline) for further processing.