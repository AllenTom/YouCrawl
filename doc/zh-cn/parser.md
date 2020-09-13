# Parser
Parser 是处理HTML文档的组件，负责将页面中的信息根据需要进行提取。

文档解析的工具为[goquery](https://github.com/PuerkitoBio/goquery)

## 自定义Parser
使用`e.AddHTMLParser(yourParser)`添加指定的处理函数即可。

## Item
`Item` 放置在`Context`中，在Parser中获取到的信息保存至此可供下一个组件(Pipeline)进行进一步处理。